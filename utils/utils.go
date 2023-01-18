package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const NotFound = "record not found"

func CreateHashPassword(unhashedPasswrod string) (string, error) {

	pass := []byte(unhashedPasswrod)

	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)

	return string(hash), err

}

func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
