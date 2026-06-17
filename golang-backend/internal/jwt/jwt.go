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
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
		"name": user.Name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return Token{}, err
	}
	return Token{Token: signedToken, ExpiresAt: expiresAt}, nil
}

func VerifyToken(tokenString string, secretKey string) (domain.UserBase, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return domain.UserBase{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user domain.UserBase
		if sub, ok := claims["sub"].(string); !ok {
			return domain.UserBase{}, errors.New("sub claim is missing or invalid")
		} else {
			user.ID = sub
		}
		if name, ok := claims["name"].(string); !ok {
			return domain.UserBase{}, errors.New("name claim is missing or invalid")
		} else {
			user.Name = name
		}
		return user, nil
	}
	return domain.UserBase{}, errors.New("invalid token")
}
