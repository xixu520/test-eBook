package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	OCR      OCRConfig      `mapstructure:"ocr"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Path         string `mapstructure:"path"`
	BusyTimeout  int    `mapstructure:"busy_timeout"`
	MaxReadConns int    `mapstructure:"max_read_conns"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type UploadConfig struct {
	Dir       string `mapstructure:"dir"`
	MaxSizeMB int    `mapstructure:"max_size_mb"`
}
type StorageConfig struct {
	Type                  string `mapstructure:"type"` // local, aliyun_oss, tencent_cos, cstcloud
	LocalPath             string `mapstructure:"local_path"`
	MaxSizeMB             int    `mapstructure:"max_size_mb"`
	AliyunEndpoint        string `mapstructure:"aliyun_endpoint"`
	AliyunAccessKeyID     string `mapstructure:"aliyun_access_key_id"`
	AliyunAccessKeySecret string `mapstructure:"aliyun_access_key_secret"`
	AliyunBucket          string `mapstructure:"aliyun_bucket"`
	TencentSecretID       string `mapstructure:"tencent_secret_id"`
	TencentSecretKey      string `mapstructure:"tencent_secret_key"`
	TencentBucketURL      string `mapstructure:"tencent_bucket_url"`
	CSTCloudEndpoint      string `mapstructure:"cstcloud_endpoint"`
	CSTCloudAccessKey     string `mapstructure:"cstcloud_access_key"`
	CSTCloudSecretKey     string `mapstructure:"cstcloud_secret_key"`
	CSTCloudBucket        string `mapstructure:"cstcloud_bucket"`
}

type OCRConfig struct {
	BaiduAPIKey        string `mapstructure:"baidu_api_key"`
	BaiduSecretKey     string `mapstructure:"baidu_secret_key"`
	PaddleToken         string `mapstructure:"paddle_token"`
	PaddleModel         string `mapstructure:"paddle_model"`
	PaddleJobURL       string `mapstructure:"paddle_job_url"`
	TimeoutSeconds     int    `mapstructure:"timeout_seconds"`
	TaskTimeoutMinutes int    `mapstructure:"task_timeout_minutes"`
	UseDocOrientationClassify bool `mapstructure:"use_doc_orientation_classify"`
	UseDocUnwarping           bool `mapstructure:"use_doc_unwarping"`
	UseChartRecognition        bool `mapstructure:"use_chart_recognition"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Output string `mapstructure:"output"`
}

var GlobalConfig *Config

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetEnvPrefix("EBOOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Warning: config.yaml not found, using env only")
		} else {
			return err
		}
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	return nil
}
