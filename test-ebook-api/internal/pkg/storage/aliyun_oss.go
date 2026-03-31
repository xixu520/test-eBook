package storage

import (
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliyunOSSStorage struct {
	client *oss.Client
	bucket *oss.Bucket
}

func NewAliyunOSSStorage(endpoint, accessID, accessSecret, bucketName string) (*AliyunOSSStorage, error) {
	client, err := oss.New(endpoint, accessID, accessSecret)
	if err != nil {
		return nil, fmt.Errorf("OSS client init failed: %v", err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("OSS bucket init failed: %v", err)
	}

	return &AliyunOSSStorage{client: client, bucket: bucket}, nil
}

func (a *AliyunOSSStorage) Save(fileName string, reader io.Reader) (string, error) {
	if err := a.bucket.PutObject(fileName, reader); err != nil {
		return "", err
	}
	// 返回对象键名，后续通过 Get 由后端代理下发或直接生成 URL
	return fileName, nil
}

func (a *AliyunOSSStorage) Get(path string) (io.ReadCloser, error) {
	return a.bucket.GetObject(path)
}

func (a *AliyunOSSStorage) Delete(path string) error {
	return a.bucket.DeleteObject(path)
}

func (a *AliyunOSSStorage) Exists(path string) (bool, error) {
	return a.bucket.IsObjectExist(path)
}

func (a *AliyunOSSStorage) TestConnection() error {
	// 尝试列出 bucket 中的对象（最多1个）来验证连通性和权限
	_, err := a.bucket.ListObjects(oss.MaxKeys(1))
	if err != nil {
		return fmt.Errorf("OSS 连接失败: %v", err)
	}
	return nil
}
