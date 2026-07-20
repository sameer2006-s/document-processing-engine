package chat

import (
	"context"
	"fmt"

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