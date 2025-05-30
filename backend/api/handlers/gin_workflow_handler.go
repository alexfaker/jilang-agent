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

// GinWorkflowHandler 处理工作流相关的请求
type GinWorkflowHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewGinWorkflowHandler 创建一个新的GinWorkflowHandler实例
func NewGinWorkflowHandler(db *gorm.DB, logger *zap.Logger) *GinWorkflowHandler {
	return &GinWorkflowHandler{
		DB:     db,
		Logger: logger,
	}
}

// CreateWorkflowRequest 创建工作流请求结构
type CreateWorkflowRequest struct {
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition" binding:"required"`
	Status      string          `json:"status"`
	AgentID     *int64          `json:"agentId"` // 关联的代理ID
}

// GetWorkflows 获取用户的工作流列表
func (h *GinWorkflowHandler) GetWorkflows(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	// 获取分页参数
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			if limit > 100 {
				limit = 100 // 限制最大查询数量
			}
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 获取筛选参数
	statusFilter := c.Query("status")

	// 构建查询 - 只查询当前用户的工作流
	query := h.DB.Model(&models.Workflow{}).Where("user_id = ?", userID.(string))

	// 应用状态筛选
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var workflows []models.Workflow
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&workflows)
	if result.Error != nil {
		h.Logger.Error("获取工作流列表失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取工作流列表失败: " + result.Error.Error(),
		})
		return
	}

	// 构造响应
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"workflows": workflows,
			"pagination": gin.H{
				"total":     total,
				"page":      offset/limit + 1,
				"page_size": limit,
				"pages":     (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetWorkflow 获取单个工作流详情
func (h *GinWorkflowHandler) GetWorkflow(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	// 获取工作流ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的工作流ID",
		})
		return
	}

	// 查询工作流 - 验证所有权
	var workflow models.Workflow
	result := h.DB.Where("id = ? AND user_id = ?", id, userID.(string)).First(&workflow)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "工作流不存在",
			})
		} else {
			h.Logger.Error("获取工作流详情失败", zap.Error(result.Error), zap.Int64("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取工作流详情失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 返回工作流详情
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   workflow,
	})
}

// CreateWorkflow 创建新工作流
func (h *GinWorkflowHandler) CreateWorkflow(c *gin.Context) {
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
	var req CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求数据格式错误: " + err.Error(),
		})
		return
	}

	// 验证 JSON 定义
	var definition map[string]interface{}
	if err := json.Unmarshal(req.Definition, &definition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "工作流定义必须是有效的JSON格式",
		})
		return
	}

	// 设置默认状态
	workflowStatus := models.WorkflowStatusDraft
	if req.Status == "active" {
		workflowStatus = models.WorkflowStatusActive
	}

	// 创建工作流
	workflow := models.Workflow{
		Name:        req.Name,
		Description: req.Description,
		Definition:  req.Definition,
		Status:      workflowStatus,
		UserID:      userID.(string),
		AgentID:     req.AgentID,
		RunCount:    0,
	}

	// 如果是从代理购买的，设置购买时间
	if req.AgentID != nil {
		now := time.Now()
		workflow.PurchasedAt = &now
	}

	// 保存到数据库
	result := h.DB.Create(&workflow)
	if result.Error != nil {
		h.Logger.Error("创建工作流失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "创建工作流失败: " + result.Error.Error(),
		})
		return
	}

	// 返回创建的工作流
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   workflow,
	})
}

// UpdateWorkflowRequest 更新工作流请求结构
type UpdateWorkflowRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition"`
	Status      string          `json:"status"`
}

// UpdateWorkflow 更新工作流
func (h *GinWorkflowHandler) UpdateWorkflow(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	// 获取工作流ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的工作流ID",
		})
		return
	}

	// 查询工作流 - 验证所有权
	var workflow models.Workflow
	result := h.DB.Where("id = ? AND user_id = ?", id, userID.(string)).First(&workflow)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "工作流不存在",
			})
		} else {
			h.Logger.Error("获取工作流失败", zap.Error(result.Error), zap.Int64("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取工作流失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 解析请求
	var req UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求数据格式错误: " + err.Error(),
		})
		return
	}

	// 验证 JSON 定义（如果提供了）
	if len(req.Definition) > 0 {
		var definition map[string]interface{}
		if err := json.Unmarshal(req.Definition, &definition); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "工作流定义必须是有效的JSON格式",
			})
			return
		}
	}

	// 准备更新数据
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if len(req.Definition) > 0 {
		updates["definition"] = req.Definition
	}
	if req.Status != "" {
		// 验证状态值
		validStatuses := []models.WorkflowStatus{
			models.WorkflowStatusDraft,
			models.WorkflowStatusActive,
			models.WorkflowStatusInactive,
			models.WorkflowStatusArchived,
		}
		validStatus := false
		for _, status := range validStatuses {
			if req.Status == string(status) {
				validStatus = true
				break
			}
		}
		if !validStatus {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "无效的工作流状态",
			})
			return
		}
		updates["status"] = req.Status
	}

	// 更新工作流
	result = h.DB.Model(&workflow).Updates(updates)
	if result.Error != nil {
		h.Logger.Error("更新工作流失败", zap.Error(result.Error), zap.Int64("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新工作流失败: " + result.Error.Error(),
		})
		return
	}

	// 重新获取更新后的工作流
	h.DB.Where("id = ? AND user_id = ?", id, userID.(string)).First(&workflow)

	// 返回更新后的工作流
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   workflow,
	})
}

// DeleteWorkflow 删除工作流
func (h *GinWorkflowHandler) DeleteWorkflow(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	// 获取工作流ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的工作流ID",
		})
		return
	}

	// 查询工作流 - 验证所有权
	var workflow models.Workflow
	result := h.DB.Where("id = ? AND user_id = ?", id, userID.(string)).First(&workflow)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "工作流不存在",
			})
		} else {
			h.Logger.Error("获取工作流失败", zap.Error(result.Error), zap.Int64("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取工作流失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 使用事务删除工作流及其执行记录
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		// 删除关联的执行记录
		if err := tx.Where("workflow_id = ?", id).Delete(&models.WorkflowExecution{}).Error; err != nil {
			return err
		}

		// 删除工作流
		if err := tx.Delete(&workflow).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		h.Logger.Error("删除工作流失败", zap.Error(err), zap.Int64("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "删除工作流失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "工作流删除成功",
	})
}
