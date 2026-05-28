package utils

import "golang.org/x/crypto/bcrypt"

func IsHashed(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
