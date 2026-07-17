package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	bc "github.com/sameer2006-s/document-processing-engine/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
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
	dberr := s.db.Create(newUser).Error
	if dberr != nil {
		return nil, dberr
	}
	return newUser, nil
}

func (s *AuthService) GetUserByEmail(email string) (*User, error) {
	var u User
	err := s.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &u, nil
}

func (s *AuthService) VerifySessionToken(token string) (*User, error) {
	var u User
	err := s.db.Where("session_token = ?", token).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Session token not found
		}
		return nil, err
	}
	return &u, nil
}

func (s *AuthService) VerifyPassword(user *User, password string) bool {
	return bc.CheckPasswordHash(password, user.Password)
}

func (s *AuthService) GenerateSessionToken(user *User) (string, error) {
	token := uuid.New().String()
	err := s.db.Model(user).Update("session_token", token).Error
	if err != nil {
		return "", err
	}
	return token, nil
}