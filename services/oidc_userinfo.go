package services

import (
	"net/http"
	"strings"

	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/gin-gonic/gin"
)

// UserInfo handles the OIDC UserInfo Endpoint (GET/POST /oauth2/userinfo).
// Claims returned are gated by the access token's scope per OIDC Core §5.4.
func UserInfo(ctx *gin.Context) {
	tokenString, err := helpers.ExtractTokenFromHeaders(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.OAuthErrorResponse{
			Error:            "invalid_token",
			ErrorDescription: "Missing or invalid access token",
		})
		return
	}

	// Extract claims including scope from the access token
	claims, err := helpers.ExtractOIDCAccessTokenClaims(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.OAuthErrorResponse{
			Error:            "invalid_token",
			ErrorDescription: "Invalid or expired access token",
		})
		return
	}

	// Look up user by UUID (sub claim) — consistent with ID token sub
	var user models.User
	if err := db.GetDB().Where("id = ?", claims.Subject).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, dto.OAuthErrorResponse{
			Error:            "invalid_token",
			ErrorDescription: "User not found",
		})
		return
	}

	// The sub claim MUST always be returned (OIDC Core §5.3.2)
	response := dto.UserInfoResponse{
		Sub: user.ID.String(),
	}

	// Include claims based on the access token's scope (OIDC Core §5.4)
	scopes := strings.Split(claims.Scope, " ")

	if helpers.ContainsScope(scopes, "profile") {
		response.Name = user.FirstName + " " + user.LastName
		response.GivenName = user.FirstName
		response.FamilyName = user.LastName
		response.PreferredUsername = user.Username
		if user.UpdatedAt != nil {
			response.UpdatedAt = user.UpdatedAt.Unix()
		}
	}

	if helpers.ContainsScope(scopes, "email") {
		response.Email = user.Email
		response.EmailVerified = true
	}

	ctx.JSON(http.StatusOK, response)
}
