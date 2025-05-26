package routes

import (
	"net/http"
	"time"

	"github.com/alexfaker/jilang-agent/api/handlers"
	"github.com/alexfaker/jilang-agent/api/middleware"
	"github.com/alexfaker/jilang-agent/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitGinRoutes 初始化Gin路由
func InitGinRoutes(db *gorm.DB, logger *zap.Logger, cfg *config.Config) *gin.Engine {
	// 创建Gin引擎
	r := gin.New()

	// 使用Gin的Recovery中间件
	r.Use(gin.Recovery())

	// 使用自定义的Logger中间件
	r.Use(middleware.GinLogger(logger))

	// 使用请求ID中间件
	r.Use(middleware.LoggerRequestIDMiddleware())

	// 配置CORS
	corsConfig := cors.Config{
		AllowOrigins:     cfg.Server.Cors.AllowedOrigins,
		AllowMethods:     cfg.Server.Cors.AllowedMethods,
		AllowHeaders:     cfg.Server.Cors.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// 静态文件服务
	if cfg.Server.ServeStatic {
		r.Static("/static", cfg.Server.StaticDir)
	}

	// 创建处理程序实例
	authHandler := handlers.NewGinAuthHandler(db, logger, cfg.Auth)
	userHandler := handlers.NewGinUserHandler(db, logger)
	workflowHandler := handlers.NewGinWorkflowHandler(db, logger)
	executionHandler := handlers.NewGinExecutionHandler(db, logger)
	agentHandler := handlers.NewGinAgentHandler(db, logger)
	statsHandler := handlers.NewGinStatsHandler(db, logger)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"message": "服务运行正常",
			},
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 公开路由
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.RefreshToken)

		// 获取代理分类
		api.GET("/agent-categories", agentHandler.GetAgentCategories)

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.GinAuthMiddleware(cfg.Auth.JWTSecret))
		{
			// 用户相关
			authorized.GET("/user/profile", userHandler.GetUserProfile)
			authorized.PUT("/user/profile", userHandler.UpdateUserProfile)
			authorized.POST("/user/change-password", userHandler.ChangePassword)
			authorized.GET("/user/:id", userHandler.GetUserByID)

			// 工作流相关
			authorized.GET("/workflows", workflowHandler.GetWorkflows)
			authorized.POST("/workflows", workflowHandler.CreateWorkflow)
			authorized.GET("/workflows/:id", workflowHandler.GetWorkflow)
			authorized.PUT("/workflows/:id", workflowHandler.UpdateWorkflow)
			authorized.DELETE("/workflows/:id", workflowHandler.DeleteWorkflow)

			// 执行相关
			authorized.GET("/executions", executionHandler.GetExecutions)
			authorized.POST("/workflows/:id/execute", executionHandler.ExecuteWorkflow)
			authorized.GET("/executions/:id", executionHandler.GetExecution)
			authorized.DELETE("/executions/:id", executionHandler.DeleteExecution)

			// 代理相关
			authorized.GET("/agents", agentHandler.GetAgents)
			authorized.POST("/agents", agentHandler.CreateAgent)
			authorized.GET("/agents/:id", agentHandler.GetAgent)
			authorized.PUT("/agents/:id", agentHandler.UpdateAgent)
			authorized.DELETE("/agents/:id", agentHandler.DeleteAgent)

			// 统计相关
			authorized.GET("/stats/dashboard", statsHandler.GetDashboardStats)
			authorized.GET("/stats/workflows", statsHandler.GetWorkflowStats)
			authorized.GET("/stats/executions", statsHandler.GetExecutionStats)
		}
	}

	return r
}
