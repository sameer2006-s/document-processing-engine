package chat

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChat(chat *Chat) error {
	return r.db.Create(chat).Error
}

func (r *ChatRepository) GetChatsByDocumentID(documentID uuid.UUID) ([]*Chat, error) {
	var chats []*Chat
	err := r.db.Where("document_id = ?", documentID).Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *ChatRepository) GetChatByID(id uuid.UUID) (*Chat, error) {
	var chat Chat
	err := r.db.Where("id = ?", id).First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *ChatRepository) UpdateChat(chat *Chat) error {
	return r.db.Save(chat).Error
}