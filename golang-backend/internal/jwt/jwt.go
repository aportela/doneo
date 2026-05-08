package jwt

import (
	"errors"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

type Token struct {
	Token     string
	ExpiresAt time.Time
}

func GenerateToken(user domain.User, expiresAt time.Time, secretKey string) (Token, error) {
	role := "user"
	if user.IsSuperUser {
		role = "administrator"
	}
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
		"role": role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return Token{}, err
	}
	return Token{Token: signedToken, ExpiresAt: expiresAt}, nil
}

func VerifyToken(tokenString string, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("sub claim is missing or invalid")
		}
		return sub, nil
	}
	return "", errors.New("invalid token")
}
