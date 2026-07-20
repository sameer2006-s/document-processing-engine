package chat

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chat struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index" foreignKey:"user_id" json:"user_id"`
	DocumentID uuid.UUID `gorm:"type:uuid;not null;index" foreignKey:"document_id" json:"document_id"`
	UserMessage string `gorm:"type:text;not null" json:"user_message"`
	AssistantMessage string `gorm:"type:text;not null" json:"assistant_message"`
	SystemPrompt string `gorm:"type:text;not null" json:"system_prompt"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}