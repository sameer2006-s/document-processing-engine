package temporal

import (
	"context"

	"github.com/google/uuid"
	"github.com/sameer2006-s/document-processing-engine/internal/chat"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"github.com/sameer2006-s/document-processing-engine/internal/ocr"
)

// DocumentStatusUpdater updates document rows from activities.
// Implemented by document.DocumentService (keeps temporal from importing document).
type DocumentStatusUpdater interface {
	UpdateDocumentStatus(id uuid.UUID, status document.DocumentStatus) error
}

type DocumentActivity struct {
	ocrService *ocr.OCRService
	chatService *chat.ChatService
	documentService *document.DocumentService
}

func NewDocumentActivity(ocrService *ocr.OCRService, chatService *chat.ChatService, documentService *document.DocumentService) *DocumentActivity {
	return &DocumentActivity{
		ocrService: ocrService,
		chatService: chatService,
		documentService: documentService,
	}
}

func (a *DocumentActivity) RunOCRActivity(ctx context.Context, documentID string) (string, error) {
	id, err := uuid.Parse(documentID)
	if err != nil {
		_ = a.documentService.UpdateDocumentStatus(id, document.DocumentStatusFailed);
		return "", err
	}
	_ = a.documentService.UpdateDocumentStatus(id, document.DocumentStatusOCRProcessing);
	result, err := a.ocrService.RunOCR(ctx, id)
	if err != nil {
		_ = a.documentService.UpdateDocumentStatus(id, document.DocumentStatusFailed);
		return "", err
	}

	if err := a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusOCRDone); err != nil {
		_ = a.ocrService.UpdateDocumentStatus(id, document.DocumentStatusFailed);
		return "", err
	}

	return result, nil
}

func (a *DocumentActivity) RunThumbnailActivity(ctx context.Context, documentID string) (string, error) {
	id, err := uuid.Parse(documentID)
	if err != nil {
		return "", err
	}
	if err := a.documentService.UpdateDocumentStatus(id, document.DocumentStatusThumbnailProcessing); err != nil {
		return "", err
	}
	thumbnail, err := a.ocrService.GenerateThumbnail(id)
	if err != nil {
		_ = a.documentService.UpdateDocumentStatus(id, document.DocumentStatusFailed);
		return "", err
	}
	if err := a.documentService.UpdateDocumentStatus(id, document.DocumentStatusThumbnailDone); err != nil {
		return "", err
	}
	return thumbnail, nil
}

func (a *DocumentActivity) UpdateDocumentStatusActivity(ctx context.Context, documentID string, status string) error {
	id, err := uuid.Parse(documentID)
	if err != nil {
		return err
	}
	err = a.documentService.UpdateDocumentStatus(id, document.DocumentStatus(status))
	if err != nil {
		return err
	}
	return nil
}

func (a *DocumentActivity) RunTagActivity(ctx context.Context, documentID string) (string, error) {
	id, err := uuid.Parse(documentID)
	if err != nil {
		return "", err
	}
	if err := a.documentService.UpdateDocumentStatus(id, document.DocumentStatusTagProcessing); err != nil {
		return "", err
	}
	tags, err := a.chatService.GenerateTags(ctx, id)
	if err != nil {
		_ = a.documentService.UpdateDocumentStatus(id, document.DocumentStatusFailed);
		return "", err
	}
	if err := a.documentService.UpdateDocumentStatus(id, document.DocumentStatusTagDone); err != nil {
		return "", err
	}
	return tags, nil
}