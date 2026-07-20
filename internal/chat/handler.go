package chat

import (
	"encoding/json"
	"fmt"
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
	chatRequest, documentID, userID, ok := h.parseChatRequest(c)
	if !ok {
		return
	}

	response, err := h.service.Chat(c.Request.Context(), chatRequest.UserPrompt, documentID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, ChatResponse{Response: response})
}

func (h *ChatHandler) ChatStream(c *gin.Context) {
	chatRequest, documentID, userID, ok := h.parseChatRequest(c)
	if !ok {
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	c.Writer.Flush()

	writeSSE := func(event string, payload any) error {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if event != "" {
			if _, err := fmt.Fprintf(c.Writer, "event: %s\n", event); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}

	_, err := h.service.ChatStream(
		c.Request.Context(),
		chatRequest.UserPrompt,
		documentID,
		userID,
		func(token string) error {
			if token == "" {
				return nil
			}
			return writeSSE("", map[string]string{"token": token})
		},
	)
	if err != nil {
		_ = writeSSE("error", map[string]string{"error": "internal server error"})
		return
	}

	_ = writeSSE("done", map[string]bool{"ok": true})
}

func (h *ChatHandler) parseChatRequest(c *gin.Context) (ChatRequest, uuid.UUID, uuid.UUID, bool) {
	var chatRequest ChatRequest
	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return ChatRequest{}, uuid.Nil, uuid.Nil, false
	}
	if chatRequest.UserPrompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user prompt is required"})
		return ChatRequest{}, uuid.Nil, uuid.Nil, false
	}
	documentID, err := uuid.Parse(chatRequest.DocumentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return ChatRequest{}, uuid.Nil, uuid.Nil, false
	}
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return ChatRequest{}, uuid.Nil, uuid.Nil, false
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return ChatRequest{}, uuid.Nil, uuid.Nil, false
	}
	return chatRequest, documentID, userUUID, true
}


