package ocr

import (
	"gorm.io/gorm"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"github.com/google/uuid"
)

type OCRRepository struct {
	db *gorm.DB
	minioClient document.MinIOClient
}

func NewOCRRepository(db *gorm.DB, minioClient document.MinIOClient) *OCRRepository {
	return &OCRRepository{db: db, minioClient: minioClient}
}

func (r *OCRRepository) UpdateDocumentStatus(id uuid.UUID, status document.DocumentStatus) error {
	return r.db.Model(&document.FileMetadata{}).Where("id = ?", id).Update("status", document.DocumentStatus(status)).Error
}

func (r *OCRRepository) GetFileMetadataByID(id uuid.UUID) (*document.FileMetadata, error) {
	var file document.FileMetadata
	if err := r.db.Where("id = ?", id).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *OCRRepository) UpdateFileMetadata(fileMetadata *document.FileMetadata) error {
	return r.db.Save(fileMetadata).Error
}

func (r *OCRRepository) GetFileContent(fileMetadata *document.FileMetadata) ([]byte, error) {
	return r.minioClient.GetFile(fileMetadata)
}

func (r *OCRRepository) SaveOCRResult(documentID uuid.UUID, ocrResult string) error {
	return r.db.Model(&document.FileMetadata{}).Where("id = ?", documentID).Update("ocr_result", ocrResult).Error
}

func (r *OCRRepository) UploadThumbnail(bucketName string, key string, fileContent []byte) error {
	return r.minioClient.UploadThumbnail(bucketName, key, fileContent)
}

func (r *OCRRepository) SaveThumbnailKey(documentID uuid.UUID, key string) error {
	return r.db.Model(&document.FileMetadata{}).Where("id = ?", documentID).Update("thumbnail_key", key).Error
}
