package ocr

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"gorm.io/gorm"
)

var (
	ErrDocumentNotInOCRProcessingStatus = errors.New("document is not in OCR processing status")
	ErrDocumentNotFound = errors.New("document not found")
)

type OCRService struct {
	repository *OCRRepository
}

func NewOCRService(db *gorm.DB, minioClient document.MinIOClient) *OCRService {
	repository := NewOCRRepository(db, minioClient)
	return &OCRService{
		repository: repository,
	}
}

func (s *OCRService) RunOCR(ctx context.Context, documentID uuid.UUID) (string,error) {
	sleep := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(sleep)
	fileMetadata, err := s.repository.GetFileMetadataByID(documentID)
	if err != nil {
		return "", ErrDocumentNotFound
	}
	if fileMetadata.Status != document.DocumentStatusOCRProcessing {
		return "", ErrDocumentNotInOCRProcessingStatus
	}
	fileContent, err := s.repository.GetFileContent(fileMetadata)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}

func (s *OCRService) UpdateDocumentStatus(id uuid.UUID, status document.DocumentStatus) error {
	return s.repository.UpdateDocumentStatus(id, status)
}