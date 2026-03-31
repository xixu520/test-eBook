package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, 0755)
	}
	return &LocalStorage{basePath: basePath}
}

func (l *LocalStorage) Save(fileName string, reader io.Reader) (string, error) {
	savePath := filepath.Join(l.basePath, fileName)
	out, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, reader); err != nil {
		return "", err
	}
	return savePath, nil
}

func (l *LocalStorage) Get(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (l *LocalStorage) Delete(path string) error {
	return os.Remove(path)
}

func (l *LocalStorage) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (l *LocalStorage) TestConnection() error {
	// 检查目录是否存在且可写
	info, err := os.Stat(l.basePath)
	if err != nil {
		return fmt.Errorf("存储目录不存在: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("路径 %s 不是目录", l.basePath)
	}
	// 尝试写入一个临时文件验证可写性
	tmpPath := filepath.Join(l.basePath, ".connection_test")
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("目录不可写: %v", err)
	}
	f.Close()
	os.Remove(tmpPath)
	return nil
}
