package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sameer2006-s/document-processing-engine/internal/auth"
)

type ChatHandler struct {
	service *ChatService
}

func NewChatHandler(service *ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

type ChatRequest struct {
	DocumentID string `json:"document_id"`
	UserPrompt string `json:"user_prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var chatRequest ChatRequest
	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if chatRequest.UserPrompt == "" {
		c.JSON(http.StatusBadRequest, ChatResponse{Response: "User prompt is required"})
		return
	}
	documentID, err := uuid.Parse(chatRequest.DocumentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ChatResponse{Response: "Invalid document ID"})
		return
	}
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ChatResponse{Response: "Unauthorized"})
		return
	}
	userId := userID.(uuid.UUID)
	userPrompt := chatRequest.UserPrompt
	response, err := h.service.Chat(c.Request.Context(), userPrompt, documentID, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ChatResponse{Response: "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, ChatResponse{Response: response})
}