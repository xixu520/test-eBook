package storage

import "io"

type Storage interface {
	// Save 保存文件，返回存储路径或 URL
	Save(fileName string, reader io.Reader) (string, error)
	// Get 获取文件内容流
	Get(path string) (io.ReadCloser, error)
	// Delete 删除文件
	Delete(path string) error
	// Exists 判断文件是否存在
	Exists(path string) (bool, error)
	// TestConnection 测试存储连通性
	TestConnection() error
}

// ObjectLister 用于云端孤儿文件扫描
// 仅云存储后端实现此接口
type ObjectLister interface {
	// ListObjects 列出存储中的所有对象键名
	ListObjects() ([]string, error)
}

