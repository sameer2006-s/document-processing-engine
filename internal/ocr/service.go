package ocr

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
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
	sleep := time.Duration(5) * time.Second
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

	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)
	part, err := multipartWriter.CreateFormFile("file", fileMetadata.OriginalName)
	if err != nil {
		return "", err
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return "", err
	}
	err = multipartWriter.Close()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8090/ocr", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", errors.New("failed to OCR document")
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	ocrResult := string(bodyBytes)
	err = s.SaveOCRResult(ctx, documentID, ocrResult)
	if err != nil {
		return "", err
	}
	err = s.UpdateDocumentStatus(documentID, document.DocumentStatusDone)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (s *OCRService) SaveOCRResult(ctx context.Context, documentID uuid.UUID, ocrResult string) error {
	return s.repository.SaveOCRResult(documentID, ocrResult)
}

func (s *OCRService) UpdateDocumentStatus(id uuid.UUID, status document.DocumentStatus) error {
	return s.repository.UpdateDocumentStatus(id, status)
}