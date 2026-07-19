package ocr

import (
	"context"
)

type OCRService struct {
}

func NewOCRService() *OCRService {
	return &OCRService{
	}
}

func (s *OCRService) RunOCR(ctx context.Context, documentID string) (string,error) {
	return "success", nil
}