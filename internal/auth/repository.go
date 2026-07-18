package auth

import (
	"errors"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user *User) error {
	err:= r.db.Create(user).Error
	if err!= nil {
		return errors.New("failed to create user") 
	}
	return nil
}

func (r *AuthRepository) UpdateSessionToken(user *User, token string) error {
	err:= r.db.Model(&User{}).Where("id = ?", user.ID).Update("session_token", token).Error
	if err!= nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found") 
		}
		return errors.New("failed to update session token") 
	}
	return nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Model(&User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found") 
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserById(id uuid.UUID) (*User, error) {
	var user User
	err := r.db.Model(&User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found") 
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByToken(token string) (*User, error) {
	var user User
	err := r.db.Model(&User{}).Where("session_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found") 
		}
		return nil, err
	}
	return &user, nil
}