package handlers

import (
	"net/http"
	"time"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GinStatsHandler 处理统计相关的请求
type GinStatsHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewGinStatsHandler 创建一个新的GinStatsHandler实例
func NewGinStatsHandler(db *gorm.DB, logger *zap.Logger) *GinStatsHandler {
	return &GinStatsHandler{
		DB:     db,
		Logger: logger,
	}
}

// GinDashboardStats 仪表盘统计数据结构
type GinDashboardStats struct {
	TotalWorkflows   int64                      `json:"total_workflows"`
	TotalExecutions  int64                      `json:"total_executions"`
	SuccessRate      float64                    `json:"success_rate"`
	RecentExecutions []models.WorkflowExecution `json:"recent_executions"`
	DailyStats       []DailyExecutionStat       `json:"daily_stats"`
}

// DailyExecutionStat 每日执行统计数据结构
type DailyExecutionStat struct {
	Date      string `json:"date"`
	Count     int    `json:"count"`
	Succeeded int    `json:"succeeded"`
	Failed    int    `json:"failed"`
}

// GetDashboardStats 获取仪表盘统计数据
func (h *GinStatsHandler) GetDashboardStats(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}
	uid := userID.(int64)

	// 获取工作流总数
	var totalWorkflows int64
	if err := h.DB.Model(&models.Workflow{}).Where("user_id = ?", uid).Count(&totalWorkflows).Error; err != nil {
		h.Logger.Error("获取工作流总数失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取统计数据失败: " + err.Error(),
		})
		return
	}

	// 获取执行总数
	var totalExecutions int64
	if err := h.DB.Model(&models.WorkflowExecution{}).
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ?", uid).
		Count(&totalExecutions).Error; err != nil {
		h.Logger.Error("获取执行总数失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取统计数据失败: " + err.Error(),
		})
		return
	}

	// 获取成功执行数量
	var succeededExecutions int64
	if err := h.DB.Model(&models.WorkflowExecution{}).
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ? AND workflow_executions.status = ?", uid, "succeeded").
		Count(&succeededExecutions).Error; err != nil {
		h.Logger.Error("获取成功执行数量失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取统计数据失败: " + err.Error(),
		})
		return
	}

	// 计算成功率
	var successRate float64
	if totalExecutions > 0 {
		successRate = float64(succeededExecutions) / float64(totalExecutions) * 100
	}

	// 获取最近的执行记录
	var recentExecutions []models.WorkflowExecution
	if err := h.DB.
		Preload("Workflow").
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ?", uid).
		Order("workflow_executions.created_at DESC").
		Limit(5).
		Find(&recentExecutions).Error; err != nil {
		h.Logger.Error("获取最近执行记录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取统计数据失败: " + err.Error(),
		})
		return
	}

	// 获取过去7天的每日执行统计
	now := time.Now()
	startDate := now.AddDate(0, 0, -6) // 7天前
	endDate := now

	// 准备日期映射
	dateMap := make(map[string]*DailyExecutionStat)
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dateMap[dateStr] = &DailyExecutionStat{
			Date:      dateStr,
			Count:     0,
			Succeeded: 0,
			Failed:    0,
		}
	}

	// 查询每日执行数据
	type DailyResult struct {
		Date      string
		Count     int
		Succeeded int
	}
	var dailyResults []DailyResult

	// 使用GORM查询每日总数和成功数
	if err := h.DB.Model(&models.WorkflowExecution{}).
		Select("DATE(workflow_executions.created_at) as date, COUNT(*) as count, SUM(CASE WHEN workflow_executions.status = 'succeeded' THEN 1 ELSE 0 END) as succeeded").
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ? AND workflow_executions.created_at BETWEEN ? AND ?", uid, startDate, endDate.AddDate(0, 0, 1)).
		Group("DATE(workflow_executions.created_at)").
		Scan(&dailyResults).Error; err != nil {
		h.Logger.Error("获取每日执行统计失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取统计数据失败: " + err.Error(),
		})
		return
	}

	// 填充日期映射
	for _, result := range dailyResults {
		if stat, exists := dateMap[result.Date]; exists {
			stat.Count = result.Count
			stat.Succeeded = result.Succeeded
			stat.Failed = result.Count - result.Succeeded
		}
	}

	// 转换为有序数组
	var dailyStats []DailyExecutionStat
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		if stat, exists := dateMap[dateStr]; exists {
			dailyStats = append(dailyStats, *stat)
		}
	}

	// 组装仪表盘统计数据
	stats := GinDashboardStats{
		TotalWorkflows:   totalWorkflows,
		TotalExecutions:  totalExecutions,
		SuccessRate:      successRate,
		RecentExecutions: recentExecutions,
		DailyStats:       dailyStats,
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}

// WorkflowStats 工作流统计数据结构
type WorkflowStats struct {
	WorkflowID   int64   `json:"workflow_id"`
	WorkflowName string  `json:"workflow_name"`
	TotalRuns    int64   `json:"total_runs"`
	SuccessRuns  int64   `json:"success_runs"`
	FailureRuns  int64   `json:"failure_runs"`
	SuccessRate  float64 `json:"success_rate"`
	AvgDuration  float64 `json:"avg_duration_ms"`
}

// GetWorkflowStats 获取工作流统计数据
func (h *GinStatsHandler) GetWorkflowStats(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}
	uid := userID.(int64)

	// 查询工作流统计数据
	var stats []WorkflowStats

	// 使用GORM查询工作流统计数据
	if err := h.DB.Model(&models.Workflow{}).
		Select(`
			workflows.id as workflow_id,
			workflows.name as workflow_name,
			COUNT(workflow_executions.id) as total_runs,
			SUM(CASE WHEN workflow_executions.status = 'succeeded' THEN 1 ELSE 0 END) as success_runs,
			SUM(CASE WHEN workflow_executions.status = 'failed' THEN 1 ELSE 0 END) as failure_runs,
			CASE WHEN COUNT(workflow_executions.id) > 0 THEN 
				(SUM(CASE WHEN workflow_executions.status = 'succeeded' THEN 1 ELSE 0 END) * 100.0 / COUNT(workflow_executions.id)) 
			ELSE 0 END as success_rate,
			AVG(CASE WHEN workflow_executions.end_time IS NOT NULL AND workflow_executions.start_time IS NOT NULL THEN 
				EXTRACT(EPOCH FROM (workflow_executions.end_time - workflow_executions.start_time)) * 1000 
			ELSE NULL END) as avg_duration_ms
		`).
		Joins("LEFT JOIN workflow_executions ON workflows.id = workflow_executions.workflow_id").
		Where("workflows.user_id = ?", uid).
		Group("workflows.id, workflows.name").
		Order("total_runs DESC").
		Scan(&stats).Error; err != nil {
		h.Logger.Error("获取工作流统计失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取工作流统计失败: " + err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}

// ExecutionStats 执行统计数据结构
type ExecutionStats struct {
	TotalExecutions int64                   `json:"total_executions"`
	SuccessRate     float64                 `json:"success_rate"`
	AvgDuration     float64                 `json:"avg_duration_ms"`
	StatusCounts    map[string]int64        `json:"status_counts"`
	Timeline        []TimelineExecutionStat `json:"timeline"`
}

// TimelineExecutionStat 时间线执行统计数据结构
type TimelineExecutionStat struct {
	Date      string `json:"date"`
	Count     int64  `json:"count"`
	Succeeded int64  `json:"succeeded"`
	Failed    int64  `json:"failed"`
}

// GetExecutionStats 获取执行统计数据
func (h *GinStatsHandler) GetExecutionStats(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}
	uid := userID.(int64)

	// 获取时间范围参数
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	// 解析日期
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的开始日期格式",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的结束日期格式",
		})
		return
	}
	// 将结束日期调整为当天结束时间
	endDate = endDate.AddDate(0, 0, 1)

	// 获取执行总数
	var totalExecutions int64
	if err := h.DB.Model(&models.WorkflowExecution{}).
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ? AND workflow_executions.created_at BETWEEN ? AND ?", uid, startDate, endDate).
		Count(&totalExecutions).Error; err != nil {
		h.Logger.Error("获取执行总数失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取执行统计失败: " + err.Error(),
		})
		return
	}

	// 获取状态计数
	type StatusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []StatusCount
	if err := h.DB.Model(&models.WorkflowExecution{}).
		Select("status, COUNT(*) as count").
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ? AND workflow_executions.created_at BETWEEN ? AND ?", uid, startDate, endDate).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		h.Logger.Error("获取状态计数失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取执行统计失败: " + err.Error(),
		})
		return
	}

	// 转换状态计数为映射
	statusCountMap := make(map[string]int64)
	var succeededCount int64
	for _, sc := range statusCounts {
		statusCountMap[sc.Status] = sc.Count
		if sc.Status == "succeeded" {
			succeededCount = sc.Count
		}
	}

	// 计算成功率
	var successRate float64
	if totalExecutions > 0 {
		successRate = float64(succeededCount) / float64(totalExecutions) * 100
	}

	// 获取平均执行时间
	var avgDuration float64
	if err := h.DB.Model(&models.WorkflowExecution{}).
		Select("AVG(EXTRACT(EPOCH FROM (end_time - start_time)) * 1000) as avg_duration").
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ? AND workflow_executions.created_at BETWEEN ? AND ? AND workflow_executions.end_time IS NOT NULL AND workflow_executions.start_time IS NOT NULL", uid, startDate, endDate).
		Scan(&avgDuration).Error; err != nil {
		h.Logger.Error("获取平均执行时间失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取执行统计失败: " + err.Error(),
		})
		return
	}

	// 获取时间线数据
	var timeline []TimelineExecutionStat

	// 使用GORM查询时间线数据
	if err := h.DB.Model(&models.WorkflowExecution{}).
		Select(`
			DATE(workflow_executions.created_at) as date,
			COUNT(*) as count,
			SUM(CASE WHEN workflow_executions.status = 'succeeded' THEN 1 ELSE 0 END) as succeeded,
			SUM(CASE WHEN workflow_executions.status = 'failed' THEN 1 ELSE 0 END) as failed
		`).
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ? AND workflow_executions.created_at BETWEEN ? AND ?", uid, startDate, endDate).
		Group("DATE(workflow_executions.created_at)").
		Order("date").
		Scan(&timeline).Error; err != nil {
		h.Logger.Error("获取时间线数据失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取执行统计失败: " + err.Error(),
		})
		return
	}

	// 组装执行统计数据
	stats := ExecutionStats{
		TotalExecutions: totalExecutions,
		SuccessRate:     successRate,
		AvgDuration:     avgDuration,
		StatusCounts:    statusCountMap,
		Timeline:        timeline,
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}
