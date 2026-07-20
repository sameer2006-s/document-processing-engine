package temporal

import (
	"context"

	"github.com/google/uuid"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"github.com/sameer2006-s/document-processing-engine/internal/ocr"
	"gorm.io/gorm"
)

// DocumentStatusUpdater updates document rows from activities.
// Implemented by document.DocumentService (keeps temporal from importing document).
type DocumentStatusUpdater interface {
	UpdateDocumentStatus(id uuid.UUID, status document.DocumentStatus) error
}

type OCRActivity struct {
	ocrService *ocr.OCRService
}

func NewOCRActivity(db *gorm.DB, minioClient document.MinIOClient) *OCRActivity {
	ocrService := ocr.NewOCRService(db, minioClient)
	return &OCRActivity{
		ocrService: ocrService,
	}
}

func (a *OCRActivity) RunOCRActivity(ctx context.Context, documentID string) (string, error) {
	id, err := uuid.Parse(documentID)
	if err != nil {
		_ = a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusFailed);
		return "", err
	}
	_ = a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusOCRProcessing);
	result, err := a.ocrService.RunOCR(ctx, id)
	if err != nil {
		_ = a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusFailed);
		return "", err
	}

	if err := a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusOCRDone); err != nil {
		_ = a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusFailed);
		return "", err
	}

	return result, nil
}

func (a *OCRActivity) RunThumbnailActivity(ctx context.Context, documentID string) (string, error) {
	id, err := uuid.Parse(documentID)
	if err != nil {
		return "", err
	}
	if err := a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusThumbnailProcessing); err != nil {
		return "", err
	}
	thumbnail, err := a.ocrService.GenerateThumbnail(id)
	if err != nil {
		_ = a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusFailed);
		return "", err
	}
	if err := a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusThumbnailDone); err != nil {
		return "", err
	}
	return thumbnail, nil
}

func (a *OCRActivity) UpdateDocumentStatusActivity(ctx context.Context, documentID string, status document.DocumentStatus) error {
	id, err := uuid.Parse(documentID)
	if err != nil {
		return err
	}
	err = a.ocrService.UpdateDocumentStatus(id, status)
	if err != nil {
		return err
	}
	return nil
}