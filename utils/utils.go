package utils

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const NotFound = "record not found"

type errorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(errorResponse{
		Success: false,
		Message: message,
	})
}

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
			ExpiresAt: time.Now().Add(time.Second * 1).Unix(),
			IssuedAt:  jwt.NewNumericDate(time.Now()).Unix(),
		},
	})
	return token.SignedString([]byte(secret))

}

func ValidateUserJWTMiddleware(c *fiber.Ctx) error {
	secret := os.Getenv("JWT_SECRET")
	bearerToken := c.GetReqHeaders()["Authorization"]
	log.Println(bearerToken)
	if bearerToken == "" {
		return c.SendStatus(http.StatusUnauthorized)
	}

	split := strings.Split(bearerToken, "Bearer ")
	token := split[1]

	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		log.Println(err)
		return ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	return nil
}

func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
