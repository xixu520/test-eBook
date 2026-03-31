package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type CSTCloudStorage struct {
	client     *minio.Client
	bucketName string
}

func NewCSTCloudStorage(endpoint, accessKey, secretKey, bucketName string) (*CSTCloudStorage, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true, // Assuming CSTCloud uses HTTPS
	})
	if err != nil {
		return nil, fmt.Errorf("S3 client init failed: %v", err)
	}

	return &CSTCloudStorage{client: minioClient, bucketName: bucketName}, nil
}

func (s *CSTCloudStorage) Save(fileName string, reader io.Reader) (string, error) {
	// Upload as octet-stream
	_, err := s.client.PutObject(context.Background(), s.bucketName, fileName, reader, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func (s *CSTCloudStorage) Get(path string) (io.ReadCloser, error) {
	return s.client.GetObject(context.Background(), s.bucketName, path, minio.GetObjectOptions{})
}

func (s *CSTCloudStorage) Delete(path string) error {
	return s.client.RemoveObject(context.Background(), s.bucketName, path, minio.RemoveObjectOptions{})
}

func (s *CSTCloudStorage) Exists(path string) (bool, error) {
	_, err := s.client.StatObject(context.Background(), s.bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *CSTCloudStorage) TestConnection() error {
	// 检查 bucket 是否存在
	exists, err := s.client.BucketExists(context.Background(), s.bucketName)
	if err != nil {
		return fmt.Errorf("S3 连接失败: %v", err)
	}
	if !exists {
		return fmt.Errorf("Bucket '%s' 不存在", s.bucketName)
	}
	return nil
}
