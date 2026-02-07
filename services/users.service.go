package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fvrvz/auth-service-go/constants"
	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/fvrvz/gologger"
	"github.com/gin-gonic/gin"
)

func GetUsers(ctx *gin.Context) {
	var users []dto.UserDTO

	if err := db.GetDB().Model(&models.User{}).Select("users.*, first_name || ' ' || last_name AS full_name").Scan(&users).Error; err != nil {
		gologger.ERROR("Unable to get all users %+v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Failed to get all users",
			Description: err.Error(),
		})
		return
	}

	if users == nil {
		users = []dto.UserDTO{}
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse[[]dto.UserDTO]{
		Message: "Users fetched successfully",
		Data:    users,
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
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		gologger.ERROR("Failed to insert user: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create user"})
		return
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

	var user dto.UserDTO

	result := db.GetDB().Model(&models.User{}).Select("users.*, first_name || ' ' || last_name AS full_name").Where(models.User{Username: username}).Scan(&user)

	if err := result.Error; err != nil {
		gologger.ERROR("Failed to delete: %+v", err)
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

	ctx.JSON(http.StatusFound, dto.SuccessResponse[dto.UserDTO]{
		Data:    user,
		Message: "User fetched successfully",
	})
}
