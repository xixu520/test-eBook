package storage

import (
	"test-ebook-api/internal/config"
)

func NewStorage(cfg config.StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "local":
		return NewLocalStorage(cfg.LocalPath), nil
	case "aliyun_oss":
		return NewAliyunOSSStorage(cfg.AliyunEndpoint, cfg.AliyunAccessKeyID, cfg.AliyunAccessKeySecret, cfg.AliyunBucket)
	case "tencent_cos":
		return NewTencentCOSStorage(cfg.TencentBucketURL, cfg.TencentSecretID, cfg.TencentSecretKey)
	case "cstcloud":
		return NewCSTCloudStorage(cfg.CSTCloudEndpoint, cfg.CSTCloudAccessKey, cfg.CSTCloudSecretKey, cfg.CSTCloudBucket)
	default:
		return NewLocalStorage(cfg.LocalPath), nil
	}
}
