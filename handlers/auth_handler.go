package handlers

import (
	"net/http"
	"shambachain/models"
	"shambachain/services"

	"github.com/gin-gonic/gin"
)

// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
	var req models.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := services.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"email":    user.Email,
		},
	})
}

// LoginHandler handles user login and token generation
func LoginHandler(c *gin.Context) {
	var req models.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	token, user, err := services.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.UserName,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}

// GetProfileHandler retrieves the profile of the currently authenticated user
func GetProfileHandler(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type in context"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.UserName,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
