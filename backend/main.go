package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexfaker/jilang-agent/api/routes"
	"github.com/alexfaker/jilang-agent/config"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/alexfaker/jilang-agent/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 初始化日志
	logger := logger.NewLogger()
	defer logger.Sync()

	logger.Info("启动AI工作流平台服务...")

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 连接数据库
	db, err := database.ConnectGormDB(cfg.Database)
	if err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}
	logger.Info("数据库连接成功")

	// 自动迁移模型（可选，生产环境可能需要手动迁移）
	if cfg.Environment != "production" {
		if err := database.AutoMigrate(db); err != nil {
			logger.Warn("数据库自动迁移失败", zap.Error(err))
		} else {
			logger.Info("数据库自动迁移成功")
		}
	}

	// 初始化Gin路由
	router := routes.InitGinRoutes(db, logger, cfg)

	// 配置服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeout.Write) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.Timeout.Idle) * time.Second,
	}

	// 启动服务器
	logger.Info("服务器启动", zap.String("地址", addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("服务器启动失败", zap.Error(err))
		os.Exit(1)
	}
}
