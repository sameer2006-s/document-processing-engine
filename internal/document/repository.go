package document

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/sameer2006-s/document-processing-engine/internal/config"
	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
	minioClient *MinIOClient
}

func NewDocumentRepository(db *gorm.DB, cfg config.MinioConfig) (*DocumentRepository, error) {
	minioClient, err:= NewMinIOClient(cfg)
	if err!= nil {
		log.Fatalf("failed to create MinIO client: %v", err)
		return nil, errors.Join(errors.New("failed to create MinIO client"), err)
	}
	exists, err := minioClient.EnsureBuckets("documents")
	if err!= nil || !exists {
		log.Fatalf("failed to ensure buckets: %v", err)
		return nil, errors.Join(errors.New("failed to ensure buckets"), err)
	}
	
	return &DocumentRepository{db: db, minioClient: minioClient}, nil
}

func (r *DocumentRepository) CreateFileMetadata(fileMetadata *FileMetadata) error {
	err:= r.db.Create(fileMetadata).Error
	if err!= nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("failed to create file metadata") 
		}
		return errors.New("failed to create file metadata") 
	}
	return nil
}

func (r *DocumentRepository) GetFileContent(id uuid.UUID) ([]byte, error) {
	fileMetadata, err:= r.GetFileMetadataByID(id)
	if err!= nil {
		return nil, errors.New("failed to get file metadata") 
	}
	fileContent, err:= r.minioClient.GetFile(fileMetadata)
	if err!= nil {
		return nil, errors.New("failed to get file content") 
	}
	return fileContent, nil
}

func (r *DocumentRepository) GetFileMetadataByID(id uuid.UUID) (*FileMetadata, error) {
	var fileMetadata FileMetadata
	err:= r.db.Model(&FileMetadata{}).Where("id = ?", id).First(&fileMetadata).Error
	if err!= nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file metadata not found") 
		}
		return nil, errors.New("failed to get file metadata") 
	}
	return &fileMetadata, nil
}

func (r *DocumentRepository) GetFileMetadataByUserID(userID uuid.UUID) ([]FileMetadata, error) {
	var fileMetadata []FileMetadata
	err:= r.db.Model(&FileMetadata{}).Where("user_id = ?", userID).Find(&fileMetadata).Error
	if err!= nil {
		return nil, errors.New("failed to get file metadata") 
	}
	return fileMetadata, nil
}

func (r *DocumentRepository) UpdateFileMetadata(fileMetadata *FileMetadata) error {
	err:= r.db.Model(&FileMetadata{}).Where("id = ?", fileMetadata.ID).Updates(fileMetadata).Error
	if err!= nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("file metadata not found") 
		}
		return errors.New("failed to update file metadata") 
	}
	return nil
}

func (r *DocumentRepository) DeleteFileMetadata(id uuid.UUID) error {
	err:= r.db.Model(&FileMetadata{}).Where("id = ?", id).Delete(&FileMetadata{}).Error
	if err!= nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("file metadata not found") 
		}
		return errors.New("failed to delete file metadata") 
	}
	return nil
}

func (r *DocumentRepository) UploadFile(fileMetadata *FileMetadata , file []byte) error {
	return errors.Join(r.CreateFileMetadata(fileMetadata), r.minioClient.UploadFile(fileMetadata, file))
}

func (r *DocumentRepository) DeleteFile(fileMetadata *FileMetadata) error {
	return errors.Join(r.DeleteFileMetadata(fileMetadata.ID), r.minioClient.DeleteFile(fileMetadata))
}

func (r *DocumentRepository) GetPendingDocuments() ([]FileMetadata, error) {
	var fileMetadata []FileMetadata
	err:= r.db.Model(&FileMetadata{}).Where("status = ?", DocumentStatusPending).Find(&fileMetadata).Error
	if err!= nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no pending documents found") 
		}
		return nil, errors.New("failed to get pending documents") 
	}
	return fileMetadata, nil
}