package routes

import (
	"net/http"

	"github.com/alexfaker/jilang-agent/api/handlers"
	jilangMiddleware "github.com/alexfaker/jilang-agent/api/middleware"
	"github.com/alexfaker/jilang-agent/config"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// InitRoutes 初始化路由
func InitRoutes(db *database.DB, logger *zap.SugaredLogger, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// 中间件
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60))

	// CORS 配置
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Server.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Server.Cors.AllowedHeaders,
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           cfg.Server.Cors.MaxAge,
	}))

	// 静态文件服务
	fileServer := http.FileServer(http.Dir(cfg.Storage.Directory))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// API 路由
	r.Route("/api", func(r chi.Router) {
		// 公共路由
		r.Group(func(r chi.Router) {
			// 身份认证
			r.Post("/auth/login", handlers.Login(db, logger, cfg))
			r.Post("/auth/register", handlers.Register(db, logger, cfg))
			r.Post("/auth/refresh", handlers.RefreshToken(db, logger, cfg))
			r.Post("/auth/logout", handlers.Logout())

			// 健康检查
			r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})

			// 代理分类路由（公开）
			r.Get("/agents/categories", handlers.GetAgentCategories(db, logger))

			// 公开代理列表路由
			r.Get("/agents/public", func(w http.ResponseWriter, r *http.Request) {
				// 设置仅查询公开代理
				r.URL.Query().Set("is_public", "true")
				handlers.ListAgents(db, logger).ServeHTTP(w, r)
			})
		})

		// 需要认证的路由
		r.Group(func(r chi.Router) {
			// JWT 认证中间件
			r.Use(jilangMiddleware.AuthMiddleware(cfg.Auth.JWTSecret))

			// 用户相关
			r.Route("/users", func(r chi.Router) {
				r.Get("/me", handlers.GetCurrentUser(db, logger))
				r.Put("/me", handlers.UpdateUser(db, logger))
				r.Put("/me/password", handlers.ChangePassword(db, logger))
			})

			// 工作流相关
			r.Route("/workflows", func(r chi.Router) {
				r.Get("/", handlers.ListWorkflows(db, logger))
				r.Post("/", handlers.CreateWorkflow(db, logger))
				r.Get("/{id}", handlers.GetWorkflow(db, logger))
				r.Put("/{id}", handlers.UpdateWorkflow(db, logger))
				r.Delete("/{id}", handlers.DeleteWorkflow(db, logger))
				r.Post("/{id}/execute", handlers.ExecuteWorkflow(db, logger))
			})

			// 工作流执行历史
			r.Route("/executions", func(r chi.Router) {
				r.Get("/", handlers.ListExecutions(db, logger))
				r.Get("/{id}", handlers.GetExecution(db, logger))
				r.Delete("/{id}/cancel", handlers.CancelExecution(db, logger))
			})

			// 代理市场
			r.Route("/agents", func(r chi.Router) {
				r.Get("/", handlers.ListAgents(db, logger))
				r.Get("/{id}", handlers.GetAgent(db, logger))
				r.Post("/", handlers.CreateAgent(db, logger))
				r.Put("/{id}", handlers.UpdateAgent(db, logger))
				r.Delete("/{id}", handlers.DeleteAgent(db, logger))
			})

			// 统计数据
			r.Route("/stats", func(r chi.Router) {
				r.Get("/dashboard", handlers.GetDashboardStats(db, logger))
				r.Get("/workflows", handlers.GetWorkflowStats(db, logger))
				r.Get("/executions", handlers.GetExecutionStats(db, logger))
			})
		})
	})

	return r
}
