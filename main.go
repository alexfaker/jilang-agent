package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexfaker/jilang-agent/api/routes"
	"github.com/alexfaker/jilang-agent/config"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/alexfaker/jilang-agent/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("无法加载配置: " + err.Error())
	}

	// 初始化日志
	logger.InitLogger(cfg.Environment == "development")
	defer logger.Sync()

	// 连接数据库
	db := database.ConnectDB(cfg.Database)
	defer db.Close()

	// 初始化路由
	router := routes.InitializeRoutes(cfg.App.JWTSecret)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	// 在单独的goroutine中启动服务器
	go func() {
		zap.L().Info("服务器启动", zap.String("地址", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 是 syscall.SIGINT
	// kill -9 是 syscall.SIGKILL，但无法被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("正在关闭服务器...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试优雅地关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Fatal("服务器强制关闭", zap.Error(err))
	}

	zap.L().Info("服务器已优雅地关闭")
}
