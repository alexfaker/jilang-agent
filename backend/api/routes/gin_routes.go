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
	purchaseHandler := handlers.NewGinPurchaseHandler(db, logger)
	rechargeHandler := handlers.NewGinRechargeHandler(db, logger)
	pointsHandler := handlers.NewGinPointsHandler(db, logger)
	settingsHandler := handlers.NewGinSettingsHandler(db, logger)

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

		// 工作流商店 - 公开的代理列表
		api.GET("/agents", agentHandler.GetAgents)                    // 获取公开代理列表
		api.GET("/agents/:id", agentHandler.GetAgent)                 // 获取代理详情
		api.GET("/agent-categories", agentHandler.GetAgentCategories) // 获取代理分类

		// 充值套餐 - 公开的充值选项
		api.GET("/recharge/packages", rechargeHandler.GetRechargePackages)

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.GinAuthMiddleware(cfg.Auth.JWTSecret))
		{
			// 用户相关
			authorized.GET("/user/profile", userHandler.GetUserProfile)
			authorized.PUT("/user/profile", userHandler.UpdateUserProfile)
			authorized.POST("/user/change-password", userHandler.ChangePassword)
			authorized.GET("/user/:id", userHandler.GetUserByID)

			// 设置相关
			authorized.GET("/settings", settingsHandler.GetSettings)
			authorized.PUT("/settings", settingsHandler.UpdateSettings)

			// 用户工作流相关 - 用户购买的工作流实例
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

			// 购买相关
			authorized.POST("/purchase/agent", purchaseHandler.PurchaseAgent)       // 购买代理
			authorized.GET("/purchase/history", purchaseHandler.GetPurchaseHistory) // 购买历史

			// 充值相关
			authorized.POST("/recharge/order", rechargeHandler.CreateRechargeOrder)        // 创建充值订单
			authorized.GET("/recharge/orders", rechargeHandler.GetRechargeOrders)          // 获取充值订单
			authorized.GET("/recharge/orders/:id", rechargeHandler.GetRechargeOrder)       // 获取充值订单详情
			authorized.POST("/recharge/payment/callback", rechargeHandler.PaymentCallback) // 支付回调

			// 点数相关
			authorized.GET("/points/balance", pointsHandler.GetPointsBalance)              // 获取点数余额
			authorized.GET("/points/transactions", pointsHandler.GetPointsTransactions)    // 获取交易历史
			authorized.GET("/points/transactions/:id", pointsHandler.GetPointsTransaction) // 获取交易详情
			authorized.GET("/points/statistics", pointsHandler.GetPointsStatistics)        // 获取统计信息

			// 代理管理（管理员或高级用户）
			authorized.POST("/admin/agents", agentHandler.CreateAgent)       // 创建代理（管理员）
			authorized.PUT("/admin/agents/:id", agentHandler.UpdateAgent)    // 更新代理（管理员）
			authorized.DELETE("/admin/agents/:id", agentHandler.DeleteAgent) // 删除代理（管理员）

			// 统计相关
			authorized.GET("/stats/dashboard", statsHandler.GetDashboardStats)
			authorized.GET("/stats/workflows", statsHandler.GetWorkflowStats)
			authorized.GET("/stats/executions", statsHandler.GetExecutionStats)
		}
	}

	return r
}
