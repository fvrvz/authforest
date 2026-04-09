package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/fvrvz/gologger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Authorize handles the OIDC Authorization Endpoint (GET /oauth2/authorize).
// Per RFC 6749 §4.1.2.1, client_id and redirect_uri are validated before any redirect.
func Authorize(ctx *gin.Context) {
	var req dto.AuthorizeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid authorization request: " + err.Error()})
		return
	}

	// Step 1: Validate client_id (must not redirect if client is unknown)
	var client models.OAuthClient
	if err := db.GetDB().Where("client_id = ?", req.ClientID).First(&client).Error; err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Unknown client_id"})
		return
	}

	// Step 2: Validate redirect_uri (must not redirect if URI is invalid)
	if !isValidRedirectURI(client.RedirectURIs, req.RedirectURI) {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid redirect_uri"})
		return
	}

	// From here, redirect_uri is validated — safe to use redirectWithError

	// Step 3: Validate response_type
	if req.ResponseType != "code" {
		redirectWithError(ctx, req.RedirectURI, "unsupported_response_type", "Only 'code' response type is supported", req.State)
		return
	}

	// Step 4: Validate openid scope (OIDC Core §3.1.2.2)
	if req.Scope == "" {
		req.Scope = "openid"
	}
	if !helpers.ContainsScope(strings.Split(req.Scope, " "), "openid") {
		redirectWithError(ctx, req.RedirectURI, "invalid_scope", "The openid scope is required", req.State)
		return
	}

	// Step 5: PKCE validation — public clients must use PKCE
	if client.ClientType == "public" && req.CodeChallenge == "" {
		redirectWithError(ctx, req.RedirectURI, "invalid_request", "PKCE code_challenge is required for public clients", req.State)
		return
	}
	if req.CodeChallengeMethod != "" && req.CodeChallengeMethod != "S256" {
		redirectWithError(ctx, req.RedirectURI, "invalid_request", "Only S256 code_challenge_method is supported", req.State)
		return
	}

	// Render login page with the authorization params embedded
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"client_name":           client.ClientName,
		"response_type":         req.ResponseType,
		"client_id":             req.ClientID,
		"redirect_uri":          req.RedirectURI,
		"scope":                 req.Scope,
		"state":                 req.State,
		"nonce":                 req.Nonce,
		"code_challenge":        req.CodeChallenge,
		"code_challenge_method": req.CodeChallengeMethod,
	})
}

// HandleLogin processes the login form submission from the authorization flow.
func HandleLogin(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	clientID := ctx.PostForm("client_id")
	redirectURI := ctx.PostForm("redirect_uri")
	scope := ctx.PostForm("scope")
	state := ctx.PostForm("state")
	nonce := ctx.PostForm("nonce")
	codeChallenge := ctx.PostForm("code_challenge")
	codeChallengeMethod := ctx.PostForm("code_challenge_method")

	// Authenticate user
	user, err := authenticateUser(dto.LoginRequest{UserId: username, Password: password})
	if err != nil {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"error":                 "Invalid username or password",
			"client_name":           clientID,
			"client_id":             clientID,
			"redirect_uri":          redirectURI,
			"scope":                 scope,
			"state":                 state,
			"nonce":                 nonce,
			"code_challenge":        codeChallenge,
			"code_challenge_method": codeChallengeMethod,
		})
		return
	}

	// Generate authorization code
	code := uuid.New().String()

	authCode := models.AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		Username:            user.Username,
		RedirectURI:         redirectURI,
		Scope:               scope,
		Nonce:               nonce,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		ExpiresAt:           time.Now().Add(5 * time.Minute),
	}

	if err := db.GetDB().Create(&authCode).Error; err != nil {
		gologger.ERROR("Failed to store authorization code: %v", err)
		redirectWithError(ctx, redirectURI, "server_error", "Failed to generate authorization code", state)
		return
	}

	// Redirect back to client with code
	redirectURL := redirectURI + "?code=" + code
	if state != "" {
		redirectURL += "&state=" + state
	}

	ctx.Redirect(http.StatusFound, redirectURL)
}
