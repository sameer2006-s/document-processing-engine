package document

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sameer2006-s/document-processing-engine/internal/config"
)

type MinIOClient struct {
	client *minio.Client
	cfg config.MinioConfig
}

func NewMinIOClient(cfg config.MinioConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		return nil, errors.Join(errors.New("failed to connect to MinIO"), err)
	}

	return &MinIOClient{client: client, cfg: cfg}, nil
}

func (r *MinIOClient) EnsureBuckets(bucketNames ...string) (bool, error) {
	for _, bucketName := range bucketNames {
		exists, err:= r.client.BucketExists(context.Background(), bucketName)
		if err!= nil {
			return false, errors.Join(errors.New("failed to check if bucket exists"), err)
		}
		if !exists {
			err:= r.CreateBucket(bucketName)
			if err!= nil {
				return false, errors.Join(errors.New("failed to create bucket"), err)
			}
		}
	}
	return true, nil
}

func (r *MinIOClient) CreateBucket(bucketName string) error {
	err:= r.client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err!= nil {
		if errors.Is(err, &minio.ErrorResponse{
			Code: "409",
			Message: "Bucket already exists",
		}) {
			return nil
		}
		return err
	}
	return nil
}

func (r *MinIOClient) UploadFile(fileMetadata *FileMetadata , file []byte) error {
	_,err:= r.client.PutObject(context.Background(), fileMetadata.BucketName, fileMetadata.MinioKey, bytesReader(file), int64(len(file)), minio.PutObjectOptions{
		ContentType: fileMetadata.ContentType,
	})
	if err!= nil {
		return errors.New("failed to upload file") 
	}
	return nil
}

func (r *MinIOClient) DeleteFile(fileMetadata *FileMetadata) error {
	err:= r.client.RemoveObject(context.Background(), fileMetadata.BucketName, fileMetadata.MinioKey, minio.RemoveObjectOptions{})
	if err!= nil {
		return errors.New("failed to delete file") 
	}
	return nil
}


func (r *MinIOClient) GetFile(fileMetadata *FileMetadata) ([]byte, error) {
	object, err:= r.client.GetObject(context.Background(), fileMetadata.BucketName, fileMetadata.MinioKey, minio.GetObjectOptions{})
	if err!= nil {
		return nil, errors.New("failed to get file") 
	}
	defer object.Close()
	fileContent, err:= io.ReadAll(object)
	if err!= nil {
		return nil, errors.New("failed to read file") 
	}
	return fileContent, nil
}

func bytesReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}