package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/dto"
	"github.com/fvrvz/auth-service-go/helpers"
	"github.com/fvrvz/auth-service-go/models"
	"github.com/gin-gonic/gin"
)

func GetUsers(ctx *gin.Context) {
	var users []dto.UserDTO

	if err := db.GetDB().Model(&models.User{}).Select("first_name, last_name, email, first_name || ' ' || last_name AS full_name, dob, created_at").Scan(&users).Error; err != nil {
		log.Fatal("Unable to get all users")
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

	if err := db.GetDB().Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusConflict, dto.ErrorResponse{
			Error: "User already exists",
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to hash password",
		})
		return
	}

	dob, erre := time.Parse("2006-01-02", req.DOB)

	if erre != nil {
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse[gin.H]{
		Message: "User registered successfully",
		Data: gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func Delete(c *gin.Context) {
	username, isMatched := c.Params.Get("userId")

	if !isMatched {
		log.Fatal("userId param not found")
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "UserId : '" + username + "' not passed correctly",
		})
		return
	}

	result := db.GetDB().Where("username = ?", username).Delete(&models.User{})

	if err := result.Error; err != nil {
		log.Println("Failed to delete:", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Database error while deleting: " + username,
			Description: "Unexpected error",
		})
		return
	}

	if result.RowsAffected == 0 {
		log.Println("User not found:", username)
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
	username, isMatched := ctx.Params.Get("userId")

	if !isMatched {
		log.Fatal("userId param not found")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "UserId : '" + username + "' not passed correctly",
		})
		return
	}

	var user dto.UserDTO

	result := db.GetDB().Model(&models.User{}).Select("first_name, last_name, email, first_name || ' ' || last_name AS full_name, dob, created_at").Where("username = ?", username).Scan(&user)

	if err := result.Error; err != nil {
		log.Println("Failed to delete:", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:       "Database error while fetching: " + username,
			Description: "Unexpected error",
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{
			Error:       "User not found",
			Description: fmt.Sprintf("No user with ID %d exists", username),
		})
		return
	}

	ctx.JSON(http.StatusFound, dto.SuccessResponse[dto.UserDTO]{
		Data:    user,
		Message: "User fetched successfully",
	})

}
