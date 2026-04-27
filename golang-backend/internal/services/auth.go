package services

import (
	"context"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(ctx context.Context, user domain.User) error
	SignIn(ctx context.Context, email, password string) (string, string, error)
}

type authService struct {
	repository                 repositories.UserRepository
	secretKey                  string
	accessTokenExpirationDays  int
	refreshTokenExpirationDays int
}

func NewAuthService(repository repositories.UserRepository, secretKey string, accessTokenExpirationDays int, refreshTokenExpirationDays int) AuthService {
	return &authService{
		repository:                 repository,
		secretKey:                  secretKey,
		accessTokenExpirationDays:  accessTokenExpirationDays,
		refreshTokenExpirationDays: refreshTokenExpirationDays,
	}
}

func (s *authService) SignUp(ctx context.Context, user domain.User) error {
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = nil
	return s.repository.Add(ctx, user)
}

func (s *authService) SignIn(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repository.GetByEmailForVerifyCredentials(ctx, email, password)
	if err != nil {
		return "", "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}
	accessToken, err := s.generateJWT(user, time.Now().Add(1*time.Hour).Unix()) // expires in 1 hour
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.generateJWT(user, time.Now().Add(365*24*time.Hour).Unix()) // expires in 365 days
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *authService) generateJWT(user domain.User, expiration int64) (string, error) {
	role := "user"
	if user.IsSuperUser {
		role = "administrator"
	}
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  expiration,
		"iat":  time.Now().Unix(),
		"role": role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
