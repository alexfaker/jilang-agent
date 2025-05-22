package handlers

import (
	"net/http"
	"time"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/alexfaker/jilang-agent/utils"
	"go.uber.org/zap"
)

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	TotalWorkflows     int                         `json:"totalWorkflows"`
	TotalExecutions    int                         `json:"totalExecutions"`
	SuccessRate        float64                     `json:"successRate"`
	RecentExecutions   []*models.WorkflowExecution `json:"recentExecutions"`
	ExecutionsByStatus map[string]int              `json:"executionsByStatus"`
	ExecutionsByDay    []DailyExecutionStats       `json:"executionsByDay"`
}

// DailyExecutionStats 每日执行统计
type DailyExecutionStats struct {
	Date      string `json:"date"`      // 格式: YYYY-MM-DD
	Count     int    `json:"count"`     // 总数
	Succeeded int    `json:"succeeded"` // 成功数
	Failed    int    `json:"failed"`    // 失败数
}

// GetDashboardStats 获取仪表盘统计数据
func GetDashboardStats(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户身份")
			return
		}

		// 获取工作流总数
		workflowCount, err := models.CountWorkflows(db.DB, userID)
		if err != nil {
			logger.Errorw("获取工作流数量失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
			return
		}

		// 获取执行总数和成功率
		executionStats, err := models.GetExecutionStats(db.DB, userID)
		if err != nil {
			logger.Errorw("获取执行统计失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
			return
		}

		// 计算成功率
		successRate := 0.0
		if executionStats.Total > 0 {
			successRate = float64(executionStats.Succeeded) / float64(executionStats.Total) * 100
		}

		// 获取按状态分组的执行数量
		executionsByStatus := map[string]int{
			"pending":   executionStats.Pending,
			"running":   executionStats.Running,
			"succeeded": executionStats.Succeeded,
			"failed":    executionStats.Failed,
			"canceled":  executionStats.Canceled,
		}

		// 获取最近的执行记录
		recentExecutions, err := models.ListExecutions(db.DB, userID, nil, nil, 5, 0)
		if err != nil {
			logger.Errorw("获取最近执行记录失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
			return
		}

		// 获取过去7天的每日执行统计
		now := time.Now()
		startDate := now.AddDate(0, 0, -6) // 7天前
		dailyStats := make([]DailyExecutionStats, 0, 7)

		for i := 0; i < 7; i++ {
			date := startDate.AddDate(0, 0, i)
			dateStr := date.Format("2006-01-02")
			nextDate := date.AddDate(0, 0, 1)

			// 获取当天的执行统计
			stats, err := models.GetExecutionStatsByDateRange(db.DB, userID, date, nextDate)
			if err != nil {
				logger.Errorw("获取每日执行统计失败", "userID", userID, "date", dateStr, "error", err)
				continue
			}

			dailyStats = append(dailyStats, DailyExecutionStats{
				Date:      dateStr,
				Count:     stats.Total,
				Succeeded: stats.Succeeded,
				Failed:    stats.Failed,
			})
		}

		// 组装仪表盘数据
		dashboardStats := DashboardStats{
			TotalWorkflows:     workflowCount,
			TotalExecutions:    executionStats.Total,
			SuccessRate:        successRate,
			RecentExecutions:   recentExecutions,
			ExecutionsByStatus: executionsByStatus,
			ExecutionsByDay:    dailyStats,
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, dashboardStats)
	}
}

// GetWorkflowStats 获取工作流统计数据
func GetWorkflowStats(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户身份")
			return
		}

		// 获取按工作流分组的统计数据
		stats, err := models.GetExecutionStatsByWorkflow(db.DB, userID)
		if err != nil {
			logger.Errorw("获取工作流统计失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, stats)
	}
}

// GetExecutionStats 获取执行统计数据
func GetExecutionStats(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户身份")
			return
		}

		// 获取执行统计数据
		stats, err := models.GetExecutionStats(db.DB, userID)
		if err != nil {
			logger.Errorw("获取执行统计失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, stats)
	}
}
