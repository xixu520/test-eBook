package service

import (
	"fmt"
	"test-ebook-api/internal/config"
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/pkg/ocr"
	"test-ebook-api/internal/pkg/storage"
	"test-ebook-api/internal/repository"
)

type SettingService struct {
	repo *repository.SettingRepository
}

func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetSettings() ([]model.SystemSetting, error) {
	return s.repo.GetAll()
}

func (s *SettingService) SaveSettings(settings []model.SystemSetting) error {
	return s.repo.BatchSave(settings)
}

func (s *SettingService) TestOCRConnection(engine string, config map[string]interface{}) error {
	if engine == "paddleocr" {
		token, _ := config["token"].(string)
		url, _ := config["url"].(string)
		model, _ := config["model"].(string)
		if url == "" {
			url = "https://paddleocr.aistudio-app.com/api/v2/ocr/jobs"
		}
		client := ocr.NewPaddleClientWithParams(token, url)
		if model != "" {
			client.SetModel(model)
		}
		return client.TestConnection()
	}
	// TODO: Add Baidu/Tesseract support if needed
	return nil
}

func (s *SettingService) TestStorageConnection(storageType string, params map[string]interface{}) error {
	// 构建临时 StorageConfig
	cfg := config.StorageConfig{Type: storageType}

	switch storageType {
	case "local":
		cfg.LocalPath, _ = params["local_path"].(string)
		if cfg.LocalPath == "" {
			cfg.LocalPath = "uploads"
		}
	case "aliyun_oss":
		cfg.AliyunEndpoint, _ = params["aliyun_endpoint"].(string)
		cfg.AliyunAccessKeyID, _ = params["aliyun_access_key_id"].(string)
		cfg.AliyunAccessKeySecret, _ = params["aliyun_access_key_secret"].(string)
		cfg.AliyunBucket, _ = params["aliyun_bucket"].(string)
		if cfg.AliyunEndpoint == "" || cfg.AliyunAccessKeyID == "" || cfg.AliyunAccessKeySecret == "" || cfg.AliyunBucket == "" {
			return fmt.Errorf("OSS 配置不完整，请填写所有必填项")
		}
	case "tencent_cos":
		cfg.TencentBucketURL, _ = params["tencent_bucket_url"].(string)
		cfg.TencentSecretID, _ = params["tencent_secret_id"].(string)
		cfg.TencentSecretKey, _ = params["tencent_secret_key"].(string)
		if cfg.TencentBucketURL == "" || cfg.TencentSecretID == "" || cfg.TencentSecretKey == "" {
			return fmt.Errorf("COS 配置不完整，请填写所有必填项")
		}
	case "cstcloud":
		cfg.CSTCloudEndpoint, _ = params["cstcloud_endpoint"].(string)
		cfg.CSTCloudAccessKey, _ = params["cstcloud_access_key"].(string)
		cfg.CSTCloudSecretKey, _ = params["cstcloud_secret_key"].(string)
		cfg.CSTCloudBucket, _ = params["cstcloud_bucket"].(string)
		if cfg.CSTCloudEndpoint == "" || cfg.CSTCloudAccessKey == "" || cfg.CSTCloudSecretKey == "" || cfg.CSTCloudBucket == "" {
			return fmt.Errorf("S3 配置不完整，请填写所有必填项")
		}
	default:
		return fmt.Errorf("不支持的存储类型: %s", storageType)
	}

	// 使用工厂创建临时存储实例并测试
	storageInst, err := storage.NewStorage(cfg)
	if err != nil {
		return fmt.Errorf("创建存储实例失败: %v", err)
	}
	return storageInst.TestConnection()
}
