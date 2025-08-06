package services

import (
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to get all users",
			"description": err.Error(),
		})
		return
	}

	if users == nil {
		users = []dto.UserDTO{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Users fetched successfully",
		"data":    users,
	})
}

func Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid Input",
			"description": err.Error(),
		})
		return
	}

	var existing models.User

	if err := db.GetDB().Where("username = ?", req.Username).First(&existing).Error; err == nil {
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

	dob, erre := time.Parse("2006-01-02", req.DOB)

	if erre != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
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

	if err := db.GetDB().Where("username = ?", username).Delete(&models.User{}).Error; err != nil {
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
