package auth

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetApiKey(headers http.Header) (string, error) {
	authheader := headers.Get("Authorization")
	if authheader == "" {
		return "", errors.New("no auth header included in request")
	}
	splitAuth := strings.Split(authheader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}
