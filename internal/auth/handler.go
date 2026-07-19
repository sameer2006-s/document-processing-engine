package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	bc "github.com/sameer2006-s/document-processing-engine/internal/utils"
)

type AuthHandler struct {
	authService *AuthService
}

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName string `json:"lastName" binding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterUserResponse struct {
	Message string `json:"message"`
	UserId uuid.UUID `json:"user_id"`
}

type LoginUserResponse struct {
	Message string `json:"message"`
	Token string `json:"token"`
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var body RegisterUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := body.Email
	password := body.Password
	firstName := body.FirstName
	lastName := body.LastName
	user, err := h.authService.CreateUser(email, password, firstName, lastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, RegisterUserResponse{Message: "User registered successfully", UserId: user.ID})
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var body LoginUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := body.Email
	password := body.Password
	user, err := h.authService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !bc.CheckPasswordHash(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := h.authService.GenerateSessionToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate session token"})
		return
	}
	c.JSON(http.StatusOK, LoginUserResponse{Message: "User logged in successfully", Token: token})
}

