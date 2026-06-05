package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"` // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Database  string `yaml:"database"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Charset   string `yaml:"charset"`
	MaxIdle   int    `yaml:"max_idle"`
	MaxOpen   int    `yaml:"max_open"`
	ParseTime bool   `yaml:"parse_time"`
	Loc       string `yaml:"loc"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int    `yaml:"expire_time"` // 过期时间（小时）
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// CozeConfig Coze AI配置
type CozeConfig struct {
	APIKey  string `yaml:"api_key"`
	BotID   string `yaml:"bot_id"`
	APIURL  string `yaml:"api_url"`
	Timeout int    `yaml:"timeout"`
}

// AppConfig 应用配置
type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
	Coze     CozeConfig     `yaml:"coze"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *AppConfig

// LoadConfig 加载配置文件
func LoadConfig(path string) (*AppConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	GlobalConfig = &config
	return GlobalConfig, nil
}

// GetConfig 获取全局配置
func GetConfig() *AppConfig {
	return GlobalConfig
}
