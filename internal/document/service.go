package document

import "github.com/google/uuid"

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

func (s *DocumentService) UploadFile(fileMetadata *FileMetadata, file []byte) error {
	return s.repository.UploadFile(fileMetadata, file)
}