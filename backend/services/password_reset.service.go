package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/fvrvz/authforest/config"
	"github.com/fvrvz/authforest/db"
	"github.com/fvrvz/authforest/dto"
	"github.com/fvrvz/authforest/helpers"
	"github.com/fvrvz/authforest/models"
	"github.com/fvrvz/gologger"
	"github.com/gin-gonic/gin"
)

const passwordResetTokenExpiry = 30 * time.Minute

func RequestPasswordReset(ctx *gin.Context) {
	var req dto.RequestPasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid input", Description: err.Error()})
		return
	}

	// Always respond with success to prevent email enumeration
	successMsg := dto.SuccessResponse[any]{Message: "If the email exists, a password reset link has been sent."}

	var user models.User
	if err := db.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		gologger.INFO("Password reset requested for non-existent email: %s", req.Email)
		ctx.JSON(http.StatusOK, successMsg)
		return
	}

	// Generate cryptographically secure token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		gologger.ERROR("Failed to generate reset token: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to process request"})
		return
	}
	token := hex.EncodeToString(tokenBytes)

	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(passwordResetTokenExpiry),
	}

	if err := db.GetDB().Create(&resetToken).Error; err != nil {
		gologger.ERROR("Failed to store reset token: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to process request"})
		return
	}

	cfg := config.GetConfig()
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", cfg.OIDC.Issuer, token)

	emailBody := fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>Hello %s,</p>
		<p>You have requested to reset your password. Click the link below to set a new password:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>This link will expire in %d minutes.</p>
		<p>If you did not request a password reset, please ignore this email.</p>
	`, user.FirstName, resetLink, int(passwordResetTokenExpiry.Minutes()))

	if err := helpers.SendEmail(&cfg.SMTP, user.Email, "Password Reset Request", emailBody); err != nil {
		gologger.ERROR("Failed to send password reset email: %v", err)
	}

	ctx.JSON(http.StatusOK, successMsg)
}

func ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid input", Description: err.Error()})
		return
	}

	var resetToken models.PasswordResetToken
	if err := db.GetDB().Where("token = ? AND used = false AND expires_at > ?", req.Token, time.Now()).First(&resetToken).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid or expired reset token"})
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		gologger.ERROR("Failed to hash new password: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to reset password"})
		return
	}

	// Update password
	if err := db.GetDB().Model(&models.User{}).Where("id = ?", resetToken.UserID).Update("password", hashedPassword).Error; err != nil {
		gologger.ERROR("Failed to update password: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to reset password"})
		return
	}

	// Mark token as used
	db.GetDB().Model(&resetToken).Update("used", true)

	// Invalidate all other unused reset tokens for this user
	db.GetDB().Model(&models.PasswordResetToken{}).Where("user_id = ? AND used = false", resetToken.UserID).Update("used", true)

	ctx.JSON(http.StatusOK, dto.SuccessResponse[any]{Message: "Password reset successfully"})
}
