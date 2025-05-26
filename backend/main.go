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
		logger.Fatalf("数据库连接失败: %v", err)
	}
	logger.Info("数据库连接成功")

	// 自动迁移模型（可选，生产环境可能需要手动迁移）
	if cfg.Environment != "production" {
		if err := db.AutoMigrate(); err != nil {
			logger.Warnf("数据库自动迁移失败: %v", err)
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
	logger.Infof("服务器启动在 %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("服务器启动失败: %v", err)
		os.Exit(1)
	}
}
