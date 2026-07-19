package temporal

import (
	"context"

	"github.com/sameer2006-s/document-processing-engine/internal/ocr"
)

type OCRActivity struct {
	ocrService *ocr.OCRService
}

func NewOCRActivity() *OCRActivity {
	return &OCRActivity{
		ocrService: ocr.NewOCRService(),
	}
}

func (s *OCRActivity) RunOCRActivity(ctx context.Context, documentID string) (string,error) {
	return s.ocrService.RunOCR(ctx, documentID)
}