package services

import (
	"net/http"

	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/fvrvz/gologger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterClient handles dynamic client registration (POST /api/v1/oauth2/register).
func RegisterClient(ctx *gin.Context) {
	var req dto.RegisterClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.OAuthErrorResponse{
			Error:            "invalid_request",
			ErrorDescription: err.Error(),
		})
		return
	}

	clientID := uuid.New().String()
	var hashedSecret string
	var rawSecret string

	if req.ClientType == "confidential" {
		rawSecret = uuid.New().String()
		hashed, err := bcrypt.GenerateFromPassword([]byte(rawSecret), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.OAuthErrorResponse{Error: "server_error", ErrorDescription: "Failed to generate client secret"})
			return
		}
		hashedSecret = string(hashed)
	}

	scopes := req.Scopes
	if scopes == "" {
		scopes = "openid profile email"
	}
	grantTypes := req.GrantTypes
	if grantTypes == "" {
		grantTypes = "authorization_code"
	}

	client := models.OAuthClient{
		ClientID:                 clientID,
		ClientSecret:             hashedSecret,
		ClientName:               req.ClientName,
		ClientType:               req.ClientType,
		RedirectURIs:             req.RedirectURIs,
		Scopes:                   scopes,
		GrantTypes:               grantTypes,
		AccessTokenExpiryMinutes: 15,
		RefreshTokenExpiryHours:  2,
		IDTokenExpiryMinutes:     15,
	}

	if req.AccessTokenExpiryMinutes != nil && *req.AccessTokenExpiryMinutes > 0 {
		client.AccessTokenExpiryMinutes = *req.AccessTokenExpiryMinutes
	}
	if req.RefreshTokenExpiryHours != nil && *req.RefreshTokenExpiryHours > 0 {
		client.RefreshTokenExpiryHours = *req.RefreshTokenExpiryHours
	}
	if req.IDTokenExpiryMinutes != nil && *req.IDTokenExpiryMinutes > 0 {
		client.IDTokenExpiryMinutes = *req.IDTokenExpiryMinutes
	}

	if err := db.GetDB().Create(&client).Error; err != nil {
		gologger.ERROR("Failed to register client: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.OAuthErrorResponse{Error: "server_error", ErrorDescription: "Failed to register client"})
		return
	}

	response := dto.RegisterClientResponse{
		ClientID:                 clientID,
		ClientName:               client.ClientName,
		ClientType:               client.ClientType,
		RedirectURIs:             req.RedirectURIs,
		Scopes:                   scopes,
		GrantTypes:               grantTypes,
		AccessTokenExpiryMinutes: client.AccessTokenExpiryMinutes,
		RefreshTokenExpiryHours:  client.RefreshTokenExpiryHours,
		IDTokenExpiryMinutes:     client.IDTokenExpiryMinutes,
		CreatedAt:                client.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if req.ClientType == "confidential" {
		response.ClientSecret = rawSecret
	}

	ctx.JSON(http.StatusCreated, response)
}

// ListClients returns all registered OAuth clients (GET /api/v1/oauth2/clients).
func ListClients(ctx *gin.Context) {
	var clients []models.OAuthClient
	if err := db.GetDB().Order("created_at desc").Find(&clients).Error; err != nil {
		gologger.ERROR("Failed to list clients: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.OAuthErrorResponse{Error: "server_error", ErrorDescription: "Failed to list clients"})
		return
	}

	response := make([]dto.RegisterClientResponse, len(clients))
	for i, c := range clients {
		response[i] = dto.RegisterClientResponse{
			ClientID:                 c.ClientID,
			ClientName:               c.ClientName,
			ClientType:               c.ClientType,
			RedirectURIs:             c.RedirectURIs,
			Scopes:                   c.Scopes,
			GrantTypes:               c.GrantTypes,
			AccessTokenExpiryMinutes: c.AccessTokenExpiryMinutes,
			RefreshTokenExpiryHours:  c.RefreshTokenExpiryHours,
			IDTokenExpiryMinutes:     c.IDTokenExpiryMinutes,
			CreatedAt:                c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[[]dto.RegisterClientResponse]{
		Message: "Clients fetched successfully",
		Data:    response,
	})
}

// GetClient returns a single OAuth client by client_id (GET /api/v1/oauth2/clients/:clientId).
func GetClient(ctx *gin.Context) {
	clientID := ctx.Param("clientId")

	var client models.OAuthClient
	if err := db.GetDB().Where("client_id = ?", clientID).First(&client).Error; err != nil {
		ctx.JSON(http.StatusNotFound, dto.OAuthErrorResponse{Error: "not_found", ErrorDescription: "Client not found"})
		return
	}

	ctx.JSON(http.StatusOK, dto.RegisterClientResponse{
		ClientID:                 client.ClientID,
		ClientName:               client.ClientName,
		ClientType:               client.ClientType,
		RedirectURIs:             client.RedirectURIs,
		Scopes:                   client.Scopes,
		GrantTypes:               client.GrantTypes,
		AccessTokenExpiryMinutes: client.AccessTokenExpiryMinutes,
		RefreshTokenExpiryHours:  client.RefreshTokenExpiryHours,
		IDTokenExpiryMinutes:     client.IDTokenExpiryMinutes,
		CreatedAt:                client.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// UpdateClient updates an OAuth client (PATCH /api/v1/oauth2/clients/:clientId).
func UpdateClient(ctx *gin.Context) {
	clientID := ctx.Param("clientId")

	var client models.OAuthClient
	if err := db.GetDB().Where("client_id = ?", clientID).First(&client).Error; err != nil {
		ctx.JSON(http.StatusNotFound, dto.OAuthErrorResponse{Error: "not_found", ErrorDescription: "Client not found"})
		return
	}

	var req dto.UpdateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.OAuthErrorResponse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	if req.ClientName != "" {
		client.ClientName = req.ClientName
	}
	if req.RedirectURIs != nil {
		client.RedirectURIs = req.RedirectURIs
	}
	if req.Scopes != "" {
		client.Scopes = req.Scopes
	}
	if req.GrantTypes != "" {
		client.GrantTypes = req.GrantTypes
	}
	if req.AccessTokenExpiryMinutes != nil && *req.AccessTokenExpiryMinutes > 0 {
		client.AccessTokenExpiryMinutes = *req.AccessTokenExpiryMinutes
	}
	if req.RefreshTokenExpiryHours != nil && *req.RefreshTokenExpiryHours > 0 {
		client.RefreshTokenExpiryHours = *req.RefreshTokenExpiryHours
	}
	if req.IDTokenExpiryMinutes != nil && *req.IDTokenExpiryMinutes > 0 {
		client.IDTokenExpiryMinutes = *req.IDTokenExpiryMinutes
	}

	if err := db.GetDB().Save(&client).Error; err != nil {
		gologger.ERROR("Failed to update client: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.OAuthErrorResponse{Error: "server_error", ErrorDescription: "Failed to update client"})
		return
	}

	ctx.JSON(http.StatusOK, dto.RegisterClientResponse{
		ClientID:                 client.ClientID,
		ClientName:               client.ClientName,
		ClientType:               client.ClientType,
		RedirectURIs:             client.RedirectURIs,
		Scopes:                   client.Scopes,
		GrantTypes:               client.GrantTypes,
		AccessTokenExpiryMinutes: client.AccessTokenExpiryMinutes,
		RefreshTokenExpiryHours:  client.RefreshTokenExpiryHours,
		IDTokenExpiryMinutes:     client.IDTokenExpiryMinutes,
		CreatedAt:                client.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// DeleteClient deletes an OAuth client (DELETE /api/v1/oauth2/clients/:clientId).
func DeleteClient(ctx *gin.Context) {
	clientID := ctx.Param("clientId")

	result := db.GetDB().Where("client_id = ?", clientID).Delete(&models.OAuthClient{})
	if result.Error != nil {
		gologger.ERROR("Failed to delete client: %v", result.Error)
		ctx.JSON(http.StatusInternalServerError, dto.OAuthErrorResponse{Error: "server_error", ErrorDescription: "Failed to delete client"})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, dto.OAuthErrorResponse{Error: "not_found", ErrorDescription: "Client not found"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[any]{Message: "Client deleted successfully"})
}

// DashboardStats returns aggregate counts for the admin dashboard (GET /api/v1/oauth2/stats).
func DashboardStats(ctx *gin.Context) {
	var stats dto.DashboardStats

	db.GetDB().Model(&models.User{}).Count(&stats.TotalUsers)
	db.GetDB().Model(&models.OAuthClient{}).Count(&stats.TotalClients)

	ctx.JSON(http.StatusOK, stats)
}
