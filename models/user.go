package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Role     string `json:"role"`
	Profile  Profile `json:"profile" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Profile struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	FarmName    string `json:"farm_name"`
	Company     string `json:"company"`
}
