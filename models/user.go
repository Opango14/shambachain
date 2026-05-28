package models

type User struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
