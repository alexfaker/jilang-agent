package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger 创建并返回新的日志记录器实例
func NewLogger() *zap.Logger {
	var config zap.Config

	// 检查环境变量确定是否为开发环境
	isDevelopment := os.Getenv("GIN_MODE") != "release"

	if isDevelopment {
		// 开发环境下使用更详细的日志
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		// 生产环境下使用结构化JSON日志
		config = zap.NewProductionConfig()
		config.OutputPaths = []string{"logs/app.log", "stdout"}
		config.ErrorOutputPaths = []string{"logs/error.log", "stderr"}

		// 确保日志目录存在
		os.MkdirAll("logs", 0755)
	}

	// 创建日志记录器
	logger, err := config.Build()
	if err != nil {
		panic("无法初始化日志记录器: " + err.Error())
	}

	// 记录初始化完成
	logger.Info("日志记录器初始化完成", zap.Bool("development", isDevelopment))

	return logger
}

// InitLogger 初始化日志记录器（保持向后兼容）
func InitLogger(isDevelopment bool) {
	var config zap.Config

	if isDevelopment {
		// 开发环境下使用更详细的日志
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		// 生产环境下使用结构化JSON日志
		config = zap.NewProductionConfig()
		config.OutputPaths = []string{"logs/app.log", "stdout"}
		config.ErrorOutputPaths = []string{"logs/error.log", "stderr"}

		// 确保日志目录存在
		os.MkdirAll("logs", 0755)
	}

	// 创建日志记录器
	logger, err := config.Build()
	if err != nil {
		panic("无法初始化日志记录器: " + err.Error())
	}

	// 替换全局logger
	zap.ReplaceGlobals(logger)

	// 记录初始化完成
	zap.L().Info("日志记录器初始化完成", zap.Bool("development", isDevelopment))
}

// Sync 刷新所有缓冲的日志
func Sync() {
	_ = zap.L().Sync()
}
