package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
)

type ChatService struct {
	provider *ChatProvider
	docService *document.DocumentService
	repository *ChatRepository
}

func NewChatService(provider *ChatProvider, docService *document.DocumentService, repository *ChatRepository) *ChatService {
	return &ChatService{provider: provider, docService: docService, repository: repository}
}

func (s *ChatService) Chat(ctx context.Context,userPrompt string, documentID uuid.UUID ,userID uuid.UUID) (string, error) {
	fileMetadata, err := s.docService.GetFileMetadata(documentID)
	if err != nil {
		return "", err
	}
	systemPrompt := fmt.Sprintf("You are a helpful assistant that can answer questions about the document. The document is: %s", fileMetadata.OCRResult)
	if fileMetadata.OCRResult == "" {
		systemPrompt = "You are a helpful assistant that can answer questions about the document."
	}
	chat := &Chat{
		UserID: userID,
		DocumentID: fileMetadata.ID,
		UserMessage: userPrompt,
		SystemPrompt: systemPrompt,
		AssistantMessage: "",
	}

	response, err := s.provider.Chat(ctx, chat.SystemPrompt, chat.UserMessage)
	if err != nil {
		return "", err
	}
	chat.AssistantMessage = response
	err = s.repository.CreateChat(chat)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (s *ChatService) ChatStream(ctx context.Context, userPrompt string, documentID uuid.UUID, userID uuid.UUID, onToken func(string) error) (full string, err error) {
	fileMetadata, err := s.docService.GetFileMetadata(documentID)
	if err != nil {
		return "", err
	}
	systemPrompt := fmt.Sprintf("You are a helpful assistant that can answer questions about the document. The document is: %s", fileMetadata.OCRResult)
	if fileMetadata.OCRResult == "" {
		systemPrompt = "You are a helpful assistant that can answer questions about the document."
	}
	chat := &Chat{
		UserID:           userID,
		DocumentID:       fileMetadata.ID,
		UserMessage:      userPrompt,
		SystemPrompt:     systemPrompt,
		AssistantMessage: "",
	}
	full, err = s.provider.ChatStream(ctx, chat.SystemPrompt, chat.UserMessage, onToken)
	if err != nil {
		return "", err
	}
	chat.AssistantMessage = full
	if err := s.repository.CreateChat(chat); err != nil {
		return full, err
	}
	return full, nil
}

func (s *ChatService) GenerateTags(ctx context.Context, documentID uuid.UUID) (string, error) {
	fileMetadata, err := s.docService.GetFileMetadata(documentID)
	if err != nil {
		return "", err
	}
	systemPrompt := fmt.Sprintf(`
	You label documents. Return ONLY valid JSON, no markdown.

	Schema: {"tags":["string",...]}
	Rules:
	- 1 to 8 tags
	- lowercase, kebab-case (e.g. invoice, id-card, contract)
	- pick from this allowlist when possible: invoice, receipt, id-card, contract, letter, form, other
	- language: use English tag names even if document is Arabic

	Document OCR:
	%s
	`, fileMetadata.OCRResult)
	if fileMetadata.OCRResult == "" {
		systemPrompt = `You label documents. Return ONLY valid JSON: {"tags":["other"]}`
	}
	response, err := s.provider.Chat(ctx, systemPrompt, "Tag this document.")
	if err != nil {
		return "", err
	}

	type tagResult struct {
		Tags []string `json:"tags"`
	}
	var out tagResult
	raw := strings.TrimSpace(response)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return "", fmt.Errorf("parse tags json: %w (raw=%q)", err, raw)
	}
	if len(out.Tags) == 0 {
		out.Tags = []string{"other"}
	}

	if err := s.docService.UpdateTags(documentID, out.Tags); err != nil {
		return "", err
	}

	encoded, err := json.Marshal(out.Tags)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
