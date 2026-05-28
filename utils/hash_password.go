package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(Password string) (password string, err error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(byte), nil
}
