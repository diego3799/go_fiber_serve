package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const NotFound = "record not found"

type JwtUser struct {
	Email string
	jwt.StandardClaims
}

func CreateHashPassword(unhashedPasswrod string) (string, error) {

	pass := []byte(unhashedPasswrod)

	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)

	return string(hash), err

}

func CreateUserJwt(email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtUser{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  jwt.NewNumericDate(time.Now()).Unix(),
		},
	})
	return token.SignedString([]byte(secret))

}

func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
