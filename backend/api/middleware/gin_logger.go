package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GinLogger 是一个Gin中间件，用于记录HTTP请求日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()
		requestID, _ := c.Get("requestID")

		// 构建日志字段
		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ip),
			zap.String("user-agent", userAgent),
			zap.Duration("latency", latency),
		}

		// 如果有请求ID，添加到日志字段
		if requestID != nil {
			fields = append(fields, zap.Any("request_id", requestID))
		}

		// 根据状态码记录不同级别的日志
		msg := fmt.Sprintf("%s %s %d %s", method, path, statusCode, latency)
		switch {
		case statusCode >= 500:
			logger.Error(msg, fields...)
		case statusCode >= 400:
			logger.Warn(msg, fields...)
		default:
			logger.Info(msg, fields...)
		}
	}
}

// LoggerRequestIDMiddleware 生成并添加请求ID到上下文
func LoggerRequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取请求ID，如果没有则生成一个新的
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 添加请求ID到上下文
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// GinRecoveryMiddleware 使用Zap记录恢复的panic
func GinRecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err interface{}) {
		// 获取请求ID
		requestID, exists := c.Get("requestID")
		requestIDStr := ""
		if exists {
			requestIDStr = requestID.(string)
		}

		// 记录panic日志
		logger.Error("Panic recovered",
			zap.Any("error", err),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("requestID", requestIDStr),
		)

		// 返回500错误
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "服务器内部错误",
		})
	})
}
