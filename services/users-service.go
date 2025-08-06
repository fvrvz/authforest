package services

import (
	"auth-service-go/helpers"
	"auth-service-go/initializers"
	"auth-service-go/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserDTO struct {
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	DOB       time.Time `json:"DOB"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func GetUsers(ctx *gin.Context) {
	var users []UserDTO

	if err := initializers.DB.Model(&models.User{}).Select("first_name, last_name, email, first_name || ' ' || last_name AS full_name, dob, created_at").Scan(&users).Error; err != nil {
		log.Fatal("Unable to get all users")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to get all users",
			"description": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Users fetched successfully",
		"data":    users,
	})
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Input",
			"details": err.Error(),
		})
		return
	}

	var existing models.User

	if err := initializers.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists",
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func Delete(c *gin.Context) {
	username, isMatched := c.Params.Get("userId")

	if !isMatched {
		log.Fatal("Username not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Username : " + username + " not found",
		})
		return
	}

	if err := initializers.DB.Where("username = ?", username).Delete(&models.User{}).Error; err != nil {
		log.Fatal("Failed to delete")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete: " + username,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": username + " Deleted successfully",
	})
}
