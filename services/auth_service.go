package services

import (
	"errors"
	"shambachain/database"
	"shambachain/models"
	"shambachain/utils"
)

// RegisterUser hashes the password and creates a new user in the database
func RegisterUser(req models.RegisterUserRequest) (*models.User, error) {
	db := database.GetDB()

	// Check if user with email or username already exists
	var existingUser models.User
	if err := db.Where("email = ? OR user_name = ?", req.Email, req.UserName).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email or username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		UserName: req.UserName,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser verifies credentials and returns a JWT token
func LoginUser(req models.LoginUserRequest) (string, error) {
	db := database.GetDB()

	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return "", errors.New("invalid email or password")
	}

	if !utils.IsHashed(user.Password, req.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID retrieves a user's details by their ID
func GetUserByID(userID uint) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
