package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("哈希加密失败:%w", err)
	}
	return string(hashedPassword), nil

}

// checkPassword
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

}
