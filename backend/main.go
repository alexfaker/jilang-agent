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

	// 连接数据库
	db, err := database.ConnectDB(cfg.Database)
	if err != nil {
		logger.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()
	logger.Info("数据库连接成功")

	// 初始化路由
	router := routes.InitRoutes(db, logger, cfg)

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
