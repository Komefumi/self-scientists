package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(plaintext string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func VerifyPassword(hashedPassword string, candidatePlaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePlaintext))
	if err != nil {
		return false
	}
	return true
}
