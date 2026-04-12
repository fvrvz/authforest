package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fvrvz/authforest/constants"
	"github.com/fvrvz/authforest/db"
	"github.com/fvrvz/authforest/dto"
	"github.com/fvrvz/authforest/helpers"
	"github.com/fvrvz/authforest/models"
	"github.com/fvrvz/gologger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetUsers(ctx *gin.Context) {
	var users []models.User

	if err := db.GetDB().Preload("Roles").Find(&users).Error; err != nil {
		gologger.ERROR("Unable to get all users %+v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Failed to get all users",
			Description: err.Error(),
		})
		return
	}

	result := make([]dto.UserDTO, len(users))
	for i, u := range users {
		result[i] = *dto.ToUserDTO(&u)
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[[]dto.UserDTO]{
		Message: "Users fetched successfully",
		Data:    result,
	})
}

func Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Invalid Input",
			Description: err.Error(),
		})
		return
	}

	var existing models.User

	if err := db.GetDB().First(&existing, models.User{Username: req.Username}).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusConflict, dto.ErrorResponse{
			Error: "User already exists",
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)

	if err != nil {
		gologger.ERROR("Failed to hash password: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to hash password",
		})
		return
	}

	dob, err := time.Parse(constants.DATE_FORMAT, req.DOB)

	if err != nil {
		gologger.ERROR("Invalid date: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid date format"})
		return
	}

	normalizedDOB := helpers.NormalizeDate(dob)

	user := models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DOB:       normalizedDOB,
		CreatedBy: req.Username,
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		gologger.ERROR("Failed to insert user: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create user"})
		return
	}

	// Assign default "user" role
	var userRole models.Role
	if err := db.GetDB().Where("name = ?", "user").First(&userRole).Error; err == nil {
		db.GetDB().Model(&user).Association("Roles").Append([]models.Role{userRole})
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse[gin.H]{
		Message: "User registered successfully",
		Data: gin.H{
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func Delete(c *gin.Context) {
	username, ok := c.Params.Get("userId")

	if !ok {
		gologger.ERROR("userId param not found")
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "UserId : '" + username + "' not passed correctly",
		})
		return
	}

	result := db.GetDB().Delete(&models.User{}, models.User{Username: username})

	if err := result.Error; err != nil {
		gologger.ERROR("Failed to delete: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Database error while deleting: " + username,
			Description: "Unexpected error",
		})
		return
	}

	if result.RowsAffected == 0 {
		gologger.WARN("User not found: %+v", username)
		c.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{
			Error:       "User not found: " + username,
			Description: "No user with that username exists",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[string]{
		Message: username + " : Deleted successfully",
	})
}

func GetUser(ctx *gin.Context) {
	username, ok := ctx.Params.Get("userId")

	if !ok {
		gologger.ERROR("userId param not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "UserId : '" + username + "' not passed correctly",
		})
		return
	}

	var user models.User

	result := db.GetDB().Preload("Roles").Where("username = ?", username).First(&user)

	if err := result.Error; err != nil {
		gologger.ERROR("Failed to fetch user: %+v", err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{
			Error:       "User not found",
			Description: fmt.Sprintf("No user with username (%s) exists", username),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[*dto.UserDTO]{
		Data:    dto.ToUserDTO(&user),
		Message: "User fetched successfully",
	})
}

func UpdateUser(ctx *gin.Context) {
	username, ok := ctx.Params.Get("userId")

	if !ok {
		gologger.ERROR("userId param not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "UserId : '" + username + "' not passed correctly",
		})
		return
	}

	var req dto.UpdateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:       "Invalid Input",
			Description: err.Error(),
		})
		return
	}

	updates := map[string]any{}

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.DOB != nil {
		dob, err := time.Parse(constants.DATE_FORMAT, *req.DOB)

		if err != nil {
			gologger.ERROR("Invalid date: %+v", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid date format"})
			return
		}

		updates["dob"] = helpers.NormalizeDate(dob)
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusNoContent, dto.ErrorResponse{
			Error: "No Fields to Update",
		})
		return
	}

	updates["updated_at"] = time.Now()
	updates["updated_by"] = username

	var user models.User

	result := db.GetDB().
		Model(&user).
		Where(models.User{Username: username}).
		Clauses(clause.Returning{}).
		Updates(&updates)

	if err := result.Error; err != nil {
		gologger.ERROR("Failed to update: %+v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Database error while fetching: " + username,
			Description: "Unexpected error",
		})
		return
	}

	if result.RowsAffected == 0 {
		gologger.INFO("User not found")
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{
			Error:       "User not found",
			Description: fmt.Sprintf("No user with username (%s) exists", username),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[dto.UserDTO]{
		Data:    *dto.ToUserDTO(&user),
		Message: "User updated successfully",
	})
}

// AdminCreateUser creates a user with roles assigned (POST /api/v1/users/create).
func AdminCreateUser(ctx *gin.Context) {
	var req dto.AdminCreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid input", Description: err.Error()})
		return
	}

	var existing models.User
	if err := db.GetDB().First(&existing, models.User{Username: req.Username}).Error; err == nil {
		ctx.JSON(http.StatusConflict, dto.ErrorResponse{Error: "User already exists"})
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		gologger.ERROR("Failed to hash password: %+v", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to hash password"})
		return
	}

	dob, err := time.Parse(constants.DATE_FORMAT, req.DOB)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid date format"})
		return
	}

	user := models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DOB:       helpers.NormalizeDate(dob),
		CreatedBy: req.Username,
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		gologger.ERROR("Failed to create user: %+v", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create user"})
		return
	}

	// Assign roles if provided, otherwise assign default "user" role
	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		if err := db.GetDB().Where("id IN ?", req.RoleIDs).Find(&roles).Error; err == nil {
			db.GetDB().Model(&user).Association("Roles").Replace(roles)
		}
	} else {
		var userRole models.Role
		if err := db.GetDB().Where("name = ?", "user").First(&userRole).Error; err == nil {
			db.GetDB().Model(&user).Association("Roles").Append([]models.Role{userRole})
		}
	}

	db.GetDB().Preload("Roles").First(&user, user.ID)

	ctx.JSON(http.StatusCreated, dto.SuccessResponse[*dto.UserDTO]{
		Message: "User created successfully",
		Data:    dto.ToUserDTO(&user),
	})
}

// AdminUpdateUser updates a user including role assignments (PATCH /api/v1/users/:userId/admin).
func AdminUpdateUser(ctx *gin.Context) {
	username, ok := ctx.Params.Get("userId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "userId param not found"})
		return
	}

	var req dto.AdminUpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid input", Description: err.Error()})
		return
	}

	var user models.User
	if err := db.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
		return
	}

	updates := map[string]any{}
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.DOB != nil {
		dob, err := time.Parse(constants.DATE_FORMAT, *req.DOB)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid date format"})
			return
		}
		updates["dob"] = helpers.NormalizeDate(dob)
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		updates["updated_by"] = username
		if err := db.GetDB().Model(&user).Updates(updates).Error; err != nil {
			gologger.ERROR("Failed to update user: %+v", err)
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update user"})
			return
		}
	}

	// Update role assignments
	if req.RoleIDs != nil {
		var roles []models.Role
		if len(req.RoleIDs) > 0 {
			db.GetDB().Where("id IN ?", req.RoleIDs).Find(&roles)
		}
		db.GetDB().Model(&user).Association("Roles").Replace(roles)
	}

	db.GetDB().Preload("Roles").First(&user, user.ID)

	ctx.JSON(http.StatusOK, dto.SuccessResponse[*dto.UserDTO]{
		Data:    dto.ToUserDTO(&user),
		Message: "User updated successfully",
	})
}
