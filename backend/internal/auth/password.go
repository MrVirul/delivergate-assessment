package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, cost int) (string, error) {
	if cost < 4 {
		cost = 12
	}
	if cost > 31 {
		cost = 12
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
