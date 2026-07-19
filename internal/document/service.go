package document

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrForbidden    = errors.New("forbidden")
)

// WorkflowStarter starts async processing for a document.
type WorkflowStarter func(documentID string) error

type DocumentService struct {
	repository    *DocumentRepository
	startWorkflow WorkflowStarter
}

func NewDocumentService(repository *DocumentRepository, startWorkflow WorkflowStarter) *DocumentService {
	return &DocumentService{repository: repository, startWorkflow: startWorkflow}
}

func (s *DocumentService) GetFileContent(id uuid.UUID) ([]byte, error) {
	return s.repository.GetFileContent(id)
}

func (s *DocumentService) CreateFileMetadata(document *FileMetadata) error {
	return s.repository.CreateFileMetadata(document)
}

func (s *DocumentService) GetFileMetadata(id uuid.UUID) (*FileMetadata, error) {
	return s.repository.GetFileMetadataByID(id)
}

func (s *DocumentService) UpdateFileMetadata(fileMetadata *FileMetadata) error {
	return s.repository.UpdateFileMetadata(fileMetadata)
}

func (s *DocumentService) UpdateDocumentStatus(id uuid.UUID, status string) error {
	return s.repository.UpdateDocumentStatus(id, DocumentStatus(status))
}

func (s *DocumentService) DeleteFileMetadata(id uuid.UUID) error {
	return s.repository.DeleteFileMetadata(id)
}

func (s *DocumentService) UploadFile(fileMetadata *FileMetadata, file []byte) (*FileMetadata, error) {
	if err := s.repository.UploadFile(fileMetadata, file); err != nil {
		return nil, err
	}

	if s.startWorkflow != nil {
		id := fileMetadata.ID.String()
		go func() {
			if err := s.startWorkflow(id); err != nil {
				log.Printf("failed to start workflow: %v", err)
			}
		}()
	}

	return fileMetadata, nil
}

func (s *DocumentService) ListDocumentsByUser(userID uuid.UUID) ([]FileMetadata, error) {
	return s.repository.GetFileMetadataByUserID(userID)
}

func (s *DocumentService) DeleteFile(id uuid.UUID, userID uuid.UUID) error {
	fileMetadata, err := s.repository.GetFileMetadataByID(id)
	if err != nil {
		return ErrFileNotFound
	}
	if fileMetadata.UserID != userID {
		return ErrForbidden
	}
	return s.repository.DeleteFile(fileMetadata)
}

func (s *DocumentService) GetPendingDocuments() ([]FileMetadata, error) {
	documents, err := s.repository.GetPendingDocuments()
	if err != nil {
		return nil, err
	}
	if len(documents) == 0 {
		return nil, nil
	}
	return documents, nil
}

func (s *DocumentService) SearchUserDocuments(userID uuid.UUID, query string) ([]FileMetadata, error) {
	documents, err := s.repository.SearchUserDocuments(userID, query)
	if err != nil {
		return nil, err
	}
	if len(documents) == 0 {
		return nil, nil
	}
	return documents, nil
}