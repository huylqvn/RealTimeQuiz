package auth

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	jwt.StandardClaims
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func GenerateJwtToken(secret string, claims jwt.Claims) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return at.SignedString([]byte(secret))
}

func ValidateJwtToken(secret, token string) (*JwtClaims, error) {
	parseToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return &JwtClaims{}, err
	}
	if !parseToken.Valid {
		return &JwtClaims{}, errors.New("invalid token")
	}

	return parseToken.Claims.(*JwtClaims), nil
}
