package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type TencentCOSStorage struct {
	client *cos.Client
}

func NewTencentCOSStorage(bucketURL, secretID, secretKey string) (*TencentCOSStorage, error) {
	u, err := url.Parse(bucketURL)
	if err != nil {
		return nil, fmt.Errorf("COS bucket URL 解析失败: %v", err)
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return &TencentCOSStorage{client: client}, nil
}

func (t *TencentCOSStorage) Save(fileName string, reader io.Reader) (string, error) {
	_, err := t.client.Object.Put(context.Background(), fileName, reader, nil)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func (t *TencentCOSStorage) Get(path string) (io.ReadCloser, error) {
	resp, err := t.client.Object.Get(context.Background(), path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (t *TencentCOSStorage) Delete(path string) error {
	_, err := t.client.Object.Delete(context.Background(), path)
	return err
}

func (t *TencentCOSStorage) Exists(path string) (bool, error) {
	return t.client.Object.IsExist(context.Background(), path)
}

func (t *TencentCOSStorage) TestConnection() error {
	// 尝试获取 bucket 信息来验证连通性
	_, _, err := t.client.Bucket.Get(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("COS 连接失败: %v", err)
	}
	return nil
}
