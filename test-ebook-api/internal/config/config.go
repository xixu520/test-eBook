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

type OCRConfig struct {
	BaiduAPIKey        string `mapstructure:"baidu_api_key"`
	BaiduSecretKey     string `mapstructure:"baidu_secret_key"`
	PaddleToken         string `mapstructure:"paddle_token"`
	PaddleModel         string `mapstructure:"paddle_model"`
	PaddleJobURL       string `mapstructure:"paddle_job_url"`
	TimeoutSeconds     int    `mapstructure:"timeout_seconds"`
	TaskTimeoutMinutes int    `mapstructure:"task_timeout_minutes"`
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
