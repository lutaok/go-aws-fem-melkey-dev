package types

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(registerUser RegisterUser) (User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)

	if err != nil {
		return User{}, fmt.Errorf("error while generating password hash - %s", err.Error())
	}

	return User{
		Username:     registerUser.Username,
		PasswordHash: string(passwordHash),
	}, nil
}

func ValidatePassword(hashedPassword string, plainTextPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
	return err
}

func CreateToken(user User) (string, error) {
	now := time.Now()
	validUntil := now.Add(time.Hour).Unix()

	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)

	// DO NOT DO THIS IN PRODUCTION: use .env or AWS Secret
	secret := "secret-test-string"

	tokenString, err := jwtToken.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
