package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const UserIDKey = "user_id"

func GetUserID(c *gin.Context) (any, error) {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return nil, errors.New("unauthorized")
	}
	if _, ok := val.(uuid.UUID); !ok {
		return nil, errors.New("invalid user id")
	}
	return val, nil
}
