package db

import (
	"github.com/sameer2006-s/document-processing-engine/internal/auth"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"github.com/sameer2006-s/document-processing-engine/internal/chat"
	"gorm.io/gorm"
)

func MyAutoMigrate(db *gorm.DB) error {
	// Drop leftover FK from the old User.Files association (type mismatch blocks AutoMigrate).
	if err := db.Exec(`ALTER TABLE IF EXISTS file_metadata DROP CONSTRAINT IF EXISTS fk_users_files`).Error; err != nil {
		return err
	}

	return db.AutoMigrate(
		&auth.User{},
		&document.FileMetadata{},
		&chat.Chat{},
	)
}
