package document

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/sameer2006-s/document-processing-engine/internal/temporal"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrForbidden    = errors.New("forbidden")
)

type DocumentService struct {
	repository *DocumentRepository
}

func NewDocumentService(repository *DocumentRepository) *DocumentService {
	return &DocumentService{repository: repository}
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

func (s *DocumentService) DeleteFileMetadata(id uuid.UUID) error {
	return s.repository.DeleteFileMetadata(id)
}

func (s *DocumentService) UploadFile(fileMetadata *FileMetadata, file []byte) (*FileMetadata, error) {
	if err := s.repository.UploadFile(fileMetadata, file); err != nil {
		return nil, err
	}

	go func() {
		cl, err := temporal.NewTemporalClient("document-processing-workflow-"+fileMetadata.ID.String(), "document-processing-task-queue")
		if err != nil {
			log.Fatal("Failed to create temporal client: ", err)
		}
		defer cl.Client.Close()
		run, runerr := cl.Client.ExecuteWorkflow(
			context.Background(),
			cl.StartWorkflowOptions,
			temporal.DocumentProcessingWorkflow,
			fileMetadata.ID.String(),
		)
		if runerr != nil {
			log.Fatal("Failed to execute workflow: ", runerr)
		}
		var result temporal.DocumentProcessingWorkflowResult
		err = run.Get(context.Background(), &result)
		if runerr != nil {
			log.Fatal("Failed to get workflow result: ", runerr)
		}
		log.Println("Workflow result: ", result)
	}()
	
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
