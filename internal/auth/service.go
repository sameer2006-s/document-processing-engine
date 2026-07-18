package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	bc "github.com/sameer2006-s/document-processing-engine/internal/utils"
)

type AuthService struct {
	authRepository *AuthRepository
}

func NewAuthService(authRepository *AuthRepository) *AuthService {
	return &AuthService{authRepository: authRepository}
}

func (s *AuthService) GenerateSessionToken(user *User) (string, error) {
	token := uuid.New().String()
	if err := s.authRepository.UpdateSessionToken(user, token); err != nil {
		return "", err
	}
	return token, nil
}
func (s *AuthService) CreateUser(email, password, firstName, lastName string) (*User, error) {
	passwordHash, hasherr := bc.HashPassword(password)
	if hasherr != nil {
		return nil, errors.New("failed to hash password")
	}
	newUser := &User{
		Email:     email,
		Password:  passwordHash,
		FirstName: firstName,
		LastName:  lastName,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	dberr := s.authRepository.CreateUser(newUser)
	if dberr != nil {
		return nil, dberr
	}
	return newUser, nil
}

func (s *AuthService) GetUserByEmail(email string) (*User, error) {
	return s.authRepository.GetUserByEmail(email)
}

func (s *AuthService) GetUserByToken(token string) (*User, error) {
	return s.authRepository.GetUserByToken(token)
}