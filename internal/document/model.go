package document

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentStatus string

const (
	DocumentStatusPending       DocumentStatus = "pending"
	DocumentStatusOCRProcessing DocumentStatus = "ocr-processing"
	DocumentStatusOCRDone DocumentStatus = "ocr-done"
	DocumentStatusThumbnailProcessing DocumentStatus = "thumbnail-processing"
	DocumentStatusThumbnailDone DocumentStatus = "thumbnail-done"
	DocumentStatusDone          DocumentStatus = "done"
	DocumentStatusFailed        DocumentStatus = "failed"
)

type FileMetadata struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	MinioKey     string         `gorm:"uniqueIndex;not null" json:"minio_key"`
	OriginalName string         `gorm:"not null" json:"original_name"`
	BucketName   string         `gorm:"not null" json:"bucket_name"`
	FileSize     int64          `json:"file_size"`
	ContentType  string         `json:"content_type"`
	Status       DocumentStatus `json:"status" gorm:"default:pending"`
	OCRResult    string         `json:"ocr_result" gorm:"type:text"`
	ThumbnailKey string         `json:"thumbnail_key"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
