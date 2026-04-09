package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/fvrvz/auth-service-go/config"
	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/fvrvz/gologger"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// TokenExchange handles the OIDC Token Endpoint (POST /oauth2/token).
// All responses include Cache-Control: no-store per OIDC Core §3.1.3.3.
func TokenExchange(ctx *gin.Context) {
	var req dto.TokenRequest
	if err := ctx.ShouldBind(&req); err != nil {
		tokenError(ctx, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	switch req.GrantType {
	case "authorization_code":
		handleAuthorizationCodeGrant(ctx, req)
	case "refresh_token":
		handleRefreshTokenGrant(ctx, req)
	default:
		tokenError(ctx, http.StatusBadRequest, "unsupported_grant_type", "Supported grant types: authorization_code, refresh_token")
	}
}

func handleAuthorizationCodeGrant(ctx *gin.Context, req dto.TokenRequest) {
	// Look up the authorization code
	var authCode models.AuthorizationCode
	if err := db.GetDB().Where("code = ?", req.Code).First(&authCode).Error; err != nil {
		tokenError(ctx, http.StatusBadRequest, "invalid_grant", "Authorization code not found")
		return
	}

	// Check if code is expired
	if time.Now().After(authCode.ExpiresAt) {
		db.GetDB().Delete(&authCode)
		tokenError(ctx, http.StatusBadRequest, "invalid_grant", "Authorization code has expired")
		return
	}

	// Check if code was already used (replay detection per OIDC Core §3.1.3.2)
	if authCode.Used {
		db.GetDB().Delete(&authCode)
		tokenError(ctx, http.StatusBadRequest, "invalid_grant", "Authorization code has already been used")
		return
	}

	// Validate client
	clientID := req.ClientID
	if clientID == "" {
		if id, _, ok := ctx.Request.BasicAuth(); ok {
			clientID = id
		}
	}

	if clientID != authCode.ClientID {
		tokenError(ctx, http.StatusBadRequest, "invalid_grant", "Client ID mismatch")
		return
	}

	var client models.OAuthClient
	if err := db.GetDB().Where("client_id = ?", clientID).First(&client).Error; err != nil {
		tokenError(ctx, http.StatusUnauthorized, "invalid_client", "Unknown client")
		return
	}

	// Validate redirect URI must match the one from the authorization request
	if req.RedirectURI != authCode.RedirectURI {
		tokenError(ctx, http.StatusBadRequest, "invalid_grant", "redirect_uri mismatch")
		return
	}

	// Authenticate confidential clients (OIDC Core §9)
	if client.ClientType == "confidential" {
		clientSecret := req.ClientSecret
		if clientSecret == "" {
			if _, secret, ok := ctx.Request.BasicAuth(); ok {
				clientSecret = secret
			}
		}
		if err := bcrypt.CompareHashAndPassword([]byte(client.ClientSecret), []byte(clientSecret)); err != nil {
			tokenError(ctx, http.StatusUnauthorized, "invalid_client", "Client authentication failed")
			return
		}
	}

	// PKCE verification (RFC 7636 §4.6)
	if authCode.CodeChallenge != "" {
		if req.CodeVerifier == "" {
			tokenError(ctx, http.StatusBadRequest, "invalid_grant", "code_verifier is required")
			return
		}
		if !verifyPKCE(authCode.CodeChallenge, req.CodeVerifier) {
			tokenError(ctx, http.StatusBadRequest, "invalid_grant", "PKCE verification failed")
			return
		}
	}

	// Mark code as used
	db.GetDB().Model(&authCode).Update("used", true)

	// Load user for claims
	var user models.User
	if err := db.GetDB().Where("username = ?", authCode.Username).First(&user).Error; err != nil {
		tokenError(ctx, http.StatusInternalServerError, "server_error", "User not found")
		return
	}

	// Generate tokens — sub claim uses user UUID for consistency across all OIDC tokens
	cfg := config.GetConfig()
	scopes := strings.Split(authCode.Scope, " ")
	authTime := authCode.CreatedAt

	accessToken, err := helpers.GenerateOIDCAccessToken(cfg, user.ID.String(), clientID, scopes, client.AccessTokenExpiryMinutes)
	if err != nil {
		gologger.ERROR("Failed to generate access token: %v", err)
		tokenError(ctx, http.StatusInternalServerError, "server_error", "Token generation failed")
		return
	}

	response := dto.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   client.AccessTokenExpiryMinutes * 60,
		Scope:       authCode.Scope,
	}

	// Generate ID token if openid scope was requested
	if helpers.ContainsScope(scopes, "openid") {
		idToken, err := helpers.GenerateIDToken(cfg, &user, clientID, authCode.Nonce, scopes, authTime, client.IDTokenExpiryMinutes)
		if err != nil {
			gologger.ERROR("Failed to generate ID token: %v", err)
			tokenError(ctx, http.StatusInternalServerError, "server_error", "ID token generation failed")
			return
		}
		response.IDToken = idToken
	}

	// Generate refresh token if offline_access scope
	if helpers.ContainsScope(scopes, "offline_access") {
		refreshToken, err := helpers.GenerateOIDCRefreshToken(cfg, user.ID.String(), clientID, scopes, client.RefreshTokenExpiryHours)
		if err != nil {
			gologger.ERROR("Failed to generate refresh token: %v", err)
			tokenError(ctx, http.StatusInternalServerError, "server_error", "Refresh token generation failed")
			return
		}
		response.RefreshToken = refreshToken

		refreshClaims, err := helpers.ExtractRSAClaims(refreshToken)
		if err != nil {
			gologger.ERROR("Failed to extract refresh token claims: %v", err)
			tokenError(ctx, http.StatusInternalServerError, "server_error", "Token generation failed")
			return
		}
		newRecord := models.AuthRefreshTokens{
			JTI:           refreshClaims.ID,
			Username:      user.Username,
			IssuedAt:      refreshClaims.IssuedAt.Time,
			ExpiresAt:     refreshClaims.ExpiresAt.Time,
			AccessTokenID: "",
			Scope:         authCode.Scope,
		}
		db.GetDB().Create(&newRecord)
	}

	ctx.Header("Cache-Control", "no-store")
	ctx.JSON(http.StatusOK, response)
}

func handleRefreshTokenGrant(ctx *gin.Context, req dto.TokenRequest) {
	if req.RefreshToken == "" {
		tokenError(ctx, http.StatusBadRequest, "invalid_request", "refresh_token is required")
		return
	}

	claims, err := helpers.ExtractRSAClaims(req.RefreshToken)
	if err != nil {
		tokenError(ctx, http.StatusBadRequest, "invalid_grant", "Invalid or expired refresh token")
		return
	}

	// Verify refresh token exists in DB
	var storedToken models.AuthRefreshTokens
	if err := db.GetDB().Where("jti = ? AND expires_at > ?", claims.ID, time.Now()).First(&storedToken).Error; err != nil {
		tokenError(ctx, http.StatusBadRequest, "invalid_grant", "Refresh token is revoked or expired")
		return
	}

	// Delete old refresh token (rotation)
	db.GetDB().Delete(&storedToken)

	// Get client ID from the request or Basic auth
	clientID := req.ClientID
	if clientID == "" {
		if id, _, ok := ctx.Request.BasicAuth(); ok {
			clientID = id
		}
	}

	// Validate audience matches
	audiences := claims.Audience
	if clientID != "" && len(audiences) > 0 && audiences[0] != clientID {
		tokenError(ctx, http.StatusBadRequest, "invalid_grant", "Client ID mismatch")
		return
	}

	// Look up user by UUID (sub claim)
	var user models.User
	if err := db.GetDB().Where("id = ?", claims.Subject).First(&user).Error; err != nil {
		tokenError(ctx, http.StatusInternalServerError, "server_error", "User not found")
		return
	}

	// Look up client for per-client token expiry settings
	var client models.OAuthClient
	if err := db.GetDB().Where("client_id = ?", clientID).First(&client).Error; err != nil {
		tokenError(ctx, http.StatusUnauthorized, "invalid_client", "Unknown client")
		return
	}

	cfg := config.GetConfig()
	scopes := strings.Split(storedToken.Scope, " ")
	if len(scopes) == 0 || scopes[0] == "" {
		scopes = []string{"openid", "profile", "email"}
	}

	accessToken, err := helpers.GenerateOIDCAccessToken(cfg, user.ID.String(), clientID, scopes, client.AccessTokenExpiryMinutes)
	if err != nil {
		tokenError(ctx, http.StatusInternalServerError, "server_error", "Token generation failed")
		return
	}

	response := dto.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   client.AccessTokenExpiryMinutes * 60,
	}

	if helpers.ContainsScope(scopes, "openid") {
		// Per OIDC Core §12.2: auth_time MUST represent the time of the original authentication
		idToken, err := helpers.GenerateIDToken(cfg, &user, clientID, "", scopes, storedToken.IssuedAt, client.IDTokenExpiryMinutes)
		if err != nil {
			gologger.ERROR("Failed to generate ID token on refresh: %v", err)
			tokenError(ctx, http.StatusInternalServerError, "server_error", "ID token generation failed")
			return
		}
		response.IDToken = idToken
	}

	// Issue new refresh token (rotation)
	newRefreshToken, err := helpers.GenerateOIDCRefreshToken(cfg, user.ID.String(), clientID, scopes, client.RefreshTokenExpiryHours)
	if err != nil {
		gologger.ERROR("Failed to generate refresh token: %v", err)
		tokenError(ctx, http.StatusInternalServerError, "server_error", "Refresh token generation failed")
		return
	}
	response.RefreshToken = newRefreshToken

	refreshClaims, err := helpers.ExtractRSAClaims(newRefreshToken)
	if err != nil {
		gologger.ERROR("Failed to extract refresh token claims: %v", err)
		tokenError(ctx, http.StatusInternalServerError, "server_error", "Token generation failed")
		return
	}
	newRecord := models.AuthRefreshTokens{
		JTI:           refreshClaims.ID,
		Username:      user.Username,
		IssuedAt:      refreshClaims.IssuedAt.Time,
		ExpiresAt:     refreshClaims.ExpiresAt.Time,
		AccessTokenID: "",
		Scope:         strings.Join(scopes, " "),
	}
	db.GetDB().Create(&newRecord)

	ctx.Header("Cache-Control", "no-store")
	ctx.JSON(http.StatusOK, response)
}

// tokenError sets Cache-Control: no-store and returns an OAuth error response per OIDC Core §3.1.3.4
func tokenError(ctx *gin.Context, status int, errorCode, description string) {
	ctx.Header("Cache-Control", "no-store")
	ctx.JSON(status, dto.OAuthErrorResponse{
		Error:            errorCode,
		ErrorDescription: description,
	})
}
