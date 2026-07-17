package db

import (
	"github.com/sameer2006-s/document-processing-engine/internal/auth"
	"gorm.io/gorm"
)

func MyAutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &auth.User{},
    )
}