package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GinExecutionHandler 处理执行相关的请求
type GinExecutionHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewGinExecutionHandler 创建一个新的GinExecutionHandler实例
func NewGinExecutionHandler(db *gorm.DB, logger *zap.Logger) *GinExecutionHandler {
	return &GinExecutionHandler{
		DB:     db,
		Logger: logger,
	}
}

// ExecuteWorkflowRequest 执行工作流请求结构
type ExecuteWorkflowRequest struct {
	WorkflowID uint            `json:"workflow_id" binding:"required"`
	Inputs     json.RawMessage `json:"inputs"`
}

// GetExecutions 获取执行记录列表
func (h *GinExecutionHandler) GetExecutions(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 获取筛选参数
	workflowID := c.Query("workflow_id")
	status := c.Query("status") // success, failed, running, canceled
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	// 构建查询
	query := h.DB.Model(&models.WorkflowExecution{})

	// 应用筛选条件
	if workflowID != "" {
		query = query.Where("workflow_id = ?", workflowID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startTime != "" {
		query = query.Where("started_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("completed_at <= ?", endTime)
	}

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var executions []models.WorkflowExecution
	result := query.Order("started_at DESC").Offset(offset).Limit(pageSize).Find(&executions)
	if result.Error != nil {
		h.Logger.Error("获取执行记录列表失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取执行记录列表失败: " + result.Error.Error(),
		})
		return
	}

	// 构造响应
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"executions": executions,
			"pagination": gin.H{
				"total":     total,
				"page":      page,
				"page_size": pageSize,
				"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetExecution 获取单个执行记录详情
func (h *GinExecutionHandler) GetExecution(c *gin.Context) {
	// 获取执行ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "执行ID不能为空",
		})
		return
	}

	// 查询执行记录
	var execution models.WorkflowExecution
	result := h.DB.First(&execution, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "执行记录不存在",
			})
		} else {
			h.Logger.Error("获取执行记录详情失败", zap.Error(result.Error), zap.String("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取执行记录详情失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 返回执行记录详情
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   execution,
	})
}

// ExecuteWorkflow 执行工作流
func (h *GinExecutionHandler) ExecuteWorkflow(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	// 解析请求
	var req ExecuteWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 查询工作流是否存在
	var workflow models.Workflow
	result := h.DB.First(&workflow, req.WorkflowID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "工作流不存在",
			})
		} else {
			h.Logger.Error("查询工作流失败", zap.Error(result.Error), zap.Uint("workflow_id", req.WorkflowID))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "查询工作流失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 检查工作流是否激活
	if workflow.Status != models.WorkflowStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "工作流未激活，无法执行",
		})
		return
	}

	// 创建执行记录
	execution := models.WorkflowExecution{
		WorkflowID: workflow.ID,
		Status:     models.ExecutionStatusRunning,
		StartedAt:  time.Now(),
		UserID:     userID.(string),
		InputData:  req.Inputs,
	}

	// 保存到数据库
	result = h.DB.Create(&execution)
	if result.Error != nil {
		h.Logger.Error("创建执行记录失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "创建执行记录失败: " + result.Error.Error(),
		})
		return
	}

	// 异步执行工作流（这里只是示例，实际实现应该使用队列或后台任务）
	go h.executeWorkflowAsync(uint(execution.ID), workflow)

	// 返回执行记录
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "工作流执行已启动",
		"data":    execution,
	})
}

// CancelExecution 取消执行
func (h *GinExecutionHandler) CancelExecution(c *gin.Context) {
	// 获取执行ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "执行ID不能为空",
		})
		return
	}

	// 查询执行记录
	var execution models.WorkflowExecution
	result := h.DB.First(&execution, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "执行记录不存在",
			})
		} else {
			h.Logger.Error("查询执行记录失败", zap.Error(result.Error), zap.String("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "查询执行记录失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 检查是否可以取消
	if execution.Status != models.ExecutionStatusRunning && execution.Status != models.ExecutionStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "只能取消处于运行中或等待中的执行",
		})
		return
	}

	// 更新执行状态
	now := time.Now()
	execution.Status = models.ExecutionStatusCancelled
	execution.CompletedAt = &now

	// 保存到数据库
	result = h.DB.Save(&execution)
	if result.Error != nil {
		h.Logger.Error("取消执行失败", zap.Error(result.Error), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "取消执行失败: " + result.Error.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "执行已取消",
		"data":    execution,
	})
}

// DeleteExecution 删除/取消执行
// 注意：此方法实际上调用的是CancelExecution方法，用于与路由配置匹配
func (h *GinExecutionHandler) DeleteExecution(c *gin.Context) {
	// 直接调用CancelExecution方法
	h.CancelExecution(c)
}

// executeWorkflowAsync 异步执行工作流
func (h *GinExecutionHandler) executeWorkflowAsync(executionID uint, workflow models.Workflow) {
	// 这里应该是实际的工作流执行逻辑
	// 在实际应用中，这可能涉及到调用外部服务、执行脚本等

	// 模拟执行过程
	time.Sleep(5 * time.Second)

	// 更新执行记录
	var execution models.WorkflowExecution
	h.DB.First(&execution, executionID)

	// 检查是否已被取消
	if execution.Status == models.ExecutionStatusCancelled {
		h.Logger.Info("执行已被取消", zap.Uint("execution_id", executionID))
		return
	}

	// 更新执行状态为成功
	execution.Status = models.ExecutionStatusSuccess
	now := time.Now()
	execution.CompletedAt = &now
	execution.OutputData = json.RawMessage(`{"result": "执行成功", "details": "这是一个模拟的执行结果"}`)

	// 保存到数据库
	result := h.DB.Save(&execution)
	if result.Error != nil {
		h.Logger.Error("更新执行记录失败", zap.Error(result.Error), zap.Uint("execution_id", executionID))
	}
}
