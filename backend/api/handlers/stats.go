package handlers

import (
	"net/http"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/alexfaker/jilang-agent/utils"
	"go.uber.org/zap"
)

// DashboardStats 仪表盘统计数据结构
type DashboardStats struct {
	TotalWorkflows     int64                      `json:"total_workflows"`
	ActiveWorkflows    int64                      `json:"active_workflows"`
	TotalExecutions    int64                      `json:"total_executions"`
	TodayExecutions    int64                      `json:"today_executions"`
	SuccessRate        float64                    `json:"success_rate"`
	RecentExecutions   []models.WorkflowExecution `json:"recent_executions"`
	ExecutionsByStatus map[string]int64           `json:"executions_by_status"`
	ExecutionsTrend    []DailyExecutions          `json:"executions_trend"`
}

// DailyExecutions 每日执行统计
type DailyExecutions struct {
	Date         string `json:"date"`
	Count        int64  `json:"count"`
	SuccessCount int64  `json:"success_count"`
	FailureCount int64  `json:"failure_count"`
}

// GetDashboardStats 获取仪表盘统计数据
func GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(int64)

	// 初始化统计数据
	stats := DashboardStats{
		ExecutionsByStatus: make(map[string]int64),
	}

	// 获取工作流总数
	totalWorkflows, err := models.CountWorkflows(database.DB, userID, "")
	if err != nil {
		zap.L().Error("获取工作流总数失败", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
		return
	}
	stats.TotalWorkflows = totalWorkflows

	// 获取活跃工作流数量
	activeWorkflows, err := models.CountWorkflows(database.DB, userID, string(models.StatusActive))
	if err != nil {
		zap.L().Error("获取活跃工作流数量失败", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
		return
	}
	stats.ActiveWorkflows = activeWorkflows

	// 获取执行总数
	totalExecutions, err := models.CountExecutions(database.DB, userID, 0, "")
	if err != nil {
		zap.L().Error("获取执行总数失败", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
		return
	}
	stats.TotalExecutions = totalExecutions

	// 获取今日执行数量
	// 注意：这里直接调用模型方法，不需要传递日期参数
	todayExecutions, err := models.CountTodayExecutions(database.DB, userID)
	if err != nil {
		zap.L().Error("获取今日执行数量失败", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
		return
	}
	stats.TodayExecutions = todayExecutions

	// 获取成功率
	if totalExecutions > 0 {
		successExecutions, err := models.CountExecutions(database.DB, userID, 0, string(models.StatusSuccess))
		if err != nil {
			zap.L().Error("获取成功执行数量失败", zap.Error(err))
			utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
			return
		}
		stats.SuccessRate = float64(successExecutions) / float64(totalExecutions) * 100
	} else {
		stats.SuccessRate = 0
	}

	// 获取最近的执行记录
	recentExecutions, err := models.ListExecutions(database.DB, userID, 0, "", 5, 0)
	if err != nil {
		zap.L().Error("获取最近执行记录失败", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取统计数据失败")
		return
	}
	stats.RecentExecutions = recentExecutions

	// 按状态统计执行数量
	for _, status := range []models.ExecutionStatus{models.StatusSuccess, models.StatusFailed, models.StatusRunning, models.StatusPending} {
		count, err := models.CountExecutions(database.DB, userID, 0, string(status))
		if err != nil {
			zap.L().Error("获取状态统计失败", zap.Error(err), zap.String("status", string(status)))
			continue
		}
		stats.ExecutionsByStatus[string(status)] = count
	}

	// 获取过去7天的执行趋势
	trend, err := models.GetExecutionTrend(database.DB, userID, 7)
	if err != nil {
		zap.L().Error("获取执行趋势失败", zap.Error(err))
	} else {
		// 转换为前端需要的格式
		for date, counts := range trend {
			dailyStats := DailyExecutions{
				Date:         date,
				Count:        counts.Total,
				SuccessCount: counts.Success,
				FailureCount: counts.Failed,
			}
			stats.ExecutionsTrend = append(stats.ExecutionsTrend, dailyStats)
		}
	}

	// 返回统计数据
	utils.RespondWithJSON(w, http.StatusOK, stats)
}
