package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 应用程序配置结构
type Config struct {
	Environment string         `json:"environment"`
	Server      ServerConfig   `json:"server"`
	Database    DatabaseConfig `json:"database"`
	Auth        AuthConfig     `json:"auth"`
	Storage     StorageConfig  `json:"storage"`
	Logging     LoggingConfig  `json:"logging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host    string        `json:"host"`
	Port    int           `json:"port"`
	Timeout TimeoutConfig `json:"timeout"`
	Cors    CorsConfig    `json:"cors"`
}

// TimeoutConfig 服务器超时配置
type TimeoutConfig struct {
	Read  int `json:"read"`  // 读取超时，秒
	Write int `json:"write"` // 写入超时，秒
	Idle  int `json:"idle"`  // 空闲超时，秒
}

// CorsConfig CORS配置
type CorsConfig struct {
	AllowedOrigins []string `json:"allowedOrigins"`
	AllowedMethods []string `json:"allowedMethods"`
	AllowedHeaders []string `json:"allowedHeaders"`
	MaxAge         int      `json:"maxAge"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslMode"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret       string `json:"jwtSecret"`
	TokenExpiration int    `json:"tokenExpiration"` // 小时
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type      string   `json:"type"` // local, s3, etc.
	Directory string   `json:"directory"`
	S3        S3Config `json:"s3"`
}

// S3Config S3存储配置
type S3Config struct {
	Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `json:"level"`
	Directory  string `json:"directory"`
	MaxSize    int    `json:"maxSize"`    // MB
	MaxBackups int    `json:"maxBackups"` // 保留的最大文件数
	MaxAge     int    `json:"maxAge"`     // 天
	Compress   bool   `json:"compress"`   // 是否压缩
}

// LoadConfig 从配置文件加载配置
func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	configPath := getConfigPath(env)
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer configFile.Close()

	var config Config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 处理环境变量覆盖
	applyEnvironmentOverrides(&config)

	// 设置默认值
	setDefaults(&config, env)

	return &config, nil
}

// getConfigPath 获取配置文件路径
func getConfigPath(env string) string {
	configName := fmt.Sprintf("config.%s.json", env)
	return filepath.Join("config", configName)
}

// applyEnvironmentOverrides 应用环境变量覆盖
func applyEnvironmentOverrides(config *Config) {
	// 数据库配置覆盖
	if os.Getenv("DB_HOST") != "" {
		config.Database.Host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		// 处理端口转换
	}
	if os.Getenv("DB_USER") != "" {
		config.Database.User = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		config.Database.Password = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_NAME") != "" {
		config.Database.Database = os.Getenv("DB_NAME")
	}

	// 服务器配置覆盖
	if os.Getenv("SERVER_PORT") != "" {
		// 处理端口转换
	}

	// JWT配置
	if os.Getenv("JWT_SECRET") != "" {
		config.Auth.JWTSecret = os.Getenv("JWT_SECRET")
	}
}

// setDefaults 设置默认配置
func setDefaults(config *Config, env string) {
	// 设置环境
	config.Environment = env

	// 服务器默认值
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Timeout.Read == 0 {
		config.Server.Timeout.Read = 15
	}
	if config.Server.Timeout.Write == 0 {
		config.Server.Timeout.Write = 15
	}
	if config.Server.Timeout.Idle == 0 {
		config.Server.Timeout.Idle = 60
	}

	// CORS默认值
	if len(config.Server.Cors.AllowedMethods) == 0 {
		config.Server.Cors.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if config.Server.Cors.MaxAge == 0 {
		config.Server.Cors.MaxAge = 300
	}

	// 数据库默认值
	if config.Database.Driver == "" {
		config.Database.Driver = "postgres"
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	// 认证默认值
	if config.Auth.TokenExpiration == 0 {
		config.Auth.TokenExpiration = 24
	}

	// 日志默认值
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Directory == "" {
		config.Logging.Directory = "logs"
	}
	if config.Logging.MaxSize == 0 {
		config.Logging.MaxSize = 100
	}
	if config.Logging.MaxBackups == 0 {
		config.Logging.MaxBackups = 7
	}
	if config.Logging.MaxAge == 0 {
		config.Logging.MaxAge = 30
	}
}
