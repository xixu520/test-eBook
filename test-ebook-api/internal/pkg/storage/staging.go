package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// StagingStorage 本地暂存存储层
// 所有上传的文件先写入此目录，再由 SyncWorker 异步同步到最终存储
type StagingStorage struct {
	basePath string
}

// NewStagingStorage 创建暂存存储，自动建目录
func NewStagingStorage(basePath string) *StagingStorage {
	if basePath == "" {
		basePath = "uploads/staging"
	}
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, 0755)
	}
	return &StagingStorage{basePath: basePath}
}

// Save 将文件写入暂存目录，返回完整暂存路径
func (s *StagingStorage) Save(fileName string, reader io.Reader) (string, error) {
	savePath := filepath.Join(s.basePath, fileName)
	out, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("创建暂存文件失败: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, reader); err != nil {
		os.Remove(savePath)
		return "", fmt.Errorf("写入暂存文件失败: %v", err)
	}
	return savePath, nil
}

// Get 读取暂存文件
func (s *StagingStorage) Get(localPath string) (io.ReadCloser, error) {
	return os.Open(localPath)
}

// Remove 删除暂存文件
func (s *StagingStorage) Remove(localPath string) error {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return nil // 文件已不存在，视为成功
	}
	return os.Remove(localPath)
}

// ListStale 列出超过指定时间的暂存文件
func (s *StagingStorage) ListStale(olderThan time.Duration) ([]string, error) {
	var staleFiles []string
	cutoff := time.Now().Add(-olderThan)

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, fmt.Errorf("读取暂存目录失败: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			staleFiles = append(staleFiles, filepath.Join(s.basePath, entry.Name()))
		}
	}
	return staleFiles, nil
}

// BasePath 返回暂存目录路径
func (s *StagingStorage) BasePath() string {
	return s.basePath
}

// Exists 检查暂存文件是否存在
func (s *StagingStorage) Exists(localPath string) bool {
	_, err := os.Stat(localPath)
	return err == nil
}
