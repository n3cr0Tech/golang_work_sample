package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(jwtKey string, username string) (string, error) {
	secretKey := []byte(jwtKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
