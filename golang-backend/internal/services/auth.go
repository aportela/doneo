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
	SignIn(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	repository repositories.UserRepository
	secretKey  string
}

func NewAuthService(repository repositories.UserRepository, secretKey string) AuthService {
	return &authService{
		repository: repository,
		secretKey:  secretKey,
	}
}

func (s *authService) SignUp(ctx context.Context, user domain.User) error {
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = nil
	return s.repository.Add(ctx, user)
}

func (s *authService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := s.repository.GetByEmailForVerifyCredentials(ctx, email, password)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}
	token, err := s.generateJWT(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) generateJWT(user domain.User) (string, error) {
	role := "user"
	if user.IsSuperUser {
		role = "administrator"
	}
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
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
