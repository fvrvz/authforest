package services

import (
	"net/http"

	"github.com/fvrvz/authforest/db"
	"github.com/fvrvz/authforest/dto"
	"github.com/fvrvz/authforest/models"
	"github.com/fvrvz/gologger"
	"github.com/gin-gonic/gin"
)

func CreateRole(ctx *gin.Context) {
	var req dto.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid input", Description: err.Error()})
		return
	}

	role := models.Role{Name: req.Name, Description: req.Description}
	if err := db.GetDB().Create(&role).Error; err != nil {
		gologger.ERROR("Failed to create role: %v", err)
		ctx.JSON(http.StatusConflict, dto.ErrorResponse{Error: "Role already exists or could not be created"})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse[dto.RoleDTO]{
		Message: "Role created successfully",
		Data: dto.RoleDTO{
			ID:          role.ID.String(),
			Name:        role.Name,
			Description: role.Description,
		},
	})
}

func ListRoles(ctx *gin.Context) {
	var roles []models.Role
	if err := db.GetDB().Order("name asc").Find(&roles).Error; err != nil {
		gologger.ERROR("Failed to list roles: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to list roles"})
		return
	}

	result := make([]dto.RoleDTO, len(roles))
	for i, r := range roles {
		result[i] = dto.RoleDTO{
			ID:          r.ID.String(),
			Name:        r.Name,
			Description: r.Description,
		}
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[[]dto.RoleDTO]{
		Message: "Roles fetched successfully",
		Data:    result,
	})
}

func GetRole(ctx *gin.Context) {
	id := ctx.Param("roleId")
	var role models.Role
	if err := db.GetDB().Where("id = ?", id).First(&role).Error; err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Role not found"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[dto.RoleDTO]{
		Message: "Role fetched successfully",
		Data: dto.RoleDTO{
			ID:          role.ID.String(),
			Name:        role.Name,
			Description: role.Description,
		},
	})
}

func UpdateRole(ctx *gin.Context) {
	id := ctx.Param("roleId")
	var role models.Role
	if err := db.GetDB().Where("id = ?", id).First(&role).Error; err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Role not found"})
		return
	}

	var req dto.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid input", Description: err.Error()})
		return
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := db.GetDB().Save(&role).Error; err != nil {
		gologger.ERROR("Failed to update role: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update role"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[dto.RoleDTO]{
		Message: "Role updated successfully",
		Data: dto.RoleDTO{
			ID:          role.ID.String(),
			Name:        role.Name,
			Description: role.Description,
		},
	})
}

func DeleteRole(ctx *gin.Context) {
	id := ctx.Param("roleId")
	result := db.GetDB().Where("id = ?", id).Delete(&models.Role{})
	if result.Error != nil {
		gologger.ERROR("Failed to delete role: %v", result.Error)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to delete role"})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Role not found"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[any]{Message: "Role deleted successfully"})
}
