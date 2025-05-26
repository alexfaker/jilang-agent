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
	Tags        []string        `json:"tags"`
}

// GetWorkflows 获取工作流列表
func (h *GinWorkflowHandler) GetWorkflows(c *gin.Context) {
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
	nameFilter := c.Query("name")
	tagFilter := c.Query("tag")
	statusFilter := c.Query("status") // active, inactive, all

	// 构建查询
	query := h.DB.Model(&models.Workflow{})

	// 应用筛选条件
	if nameFilter != "" {
		query = query.Where("name LIKE ?", "%"+nameFilter+"%")
	}
	if tagFilter != "" {
		query = query.Where("JSON_CONTAINS(tags, ?)", `"`+tagFilter+`"`)
	}
	if statusFilter == "active" {
		query = query.Where("status = ?", models.WorkflowStatusActive)
	} else if statusFilter == "inactive" {
		query = query.Where("status = ?", models.WorkflowStatusInactive)
	} else if statusFilter == "archived" {
		query = query.Where("status = ?", models.WorkflowStatusArchived)
	} else if statusFilter == "draft" {
		query = query.Where("status = ?", models.WorkflowStatusDraft)
	}

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var workflows []models.Workflow
	result := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&workflows)
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
				"page":      page,
				"page_size": pageSize,
				"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetWorkflow 获取单个工作流详情
func (h *GinWorkflowHandler) GetWorkflow(c *gin.Context) {
	// 获取工作流ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "工作流ID不能为空",
		})
		return
	}

	// 查询工作流
	var workflow models.Workflow
	result := h.DB.First(&workflow, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "工作流不存在",
			})
		} else {
			h.Logger.Error("获取工作流详情失败", zap.Error(result.Error), zap.String("id", id))
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
			"message": "无效的请求数据: " + err.Error(),
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
		UserID:      userID.(int64),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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
		"status":  "success",
		"message": "工作流创建成功",
		"data":    workflow,
	})
}

// UpdateWorkflow 更新工作流
func (h *GinWorkflowHandler) UpdateWorkflow(c *gin.Context) {
	// 获取工作流ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "工作流ID不能为空",
		})
		return
	}

	// 查询工作流是否存在
	var workflow models.Workflow
	result := h.DB.First(&workflow, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "工作流不存在",
			})
		} else {
			h.Logger.Error("查询工作流失败", zap.Error(result.Error), zap.String("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "查询工作流失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 解析请求
	var req CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 设置工作流状态
	workflowStatus := workflow.Status // 保持原有状态
	if req.Status == "active" {
		workflowStatus = models.WorkflowStatusActive
	} else if req.Status == "inactive" {
		workflowStatus = models.WorkflowStatusInactive
	} else if req.Status == "archived" {
		workflowStatus = models.WorkflowStatusArchived
	} else if req.Status == "draft" {
		workflowStatus = models.WorkflowStatusDraft
	}

	// 更新工作流
	workflow.Name = req.Name
	workflow.Description = req.Description
	workflow.Definition = req.Definition
	workflow.Status = workflowStatus
	workflow.UpdatedAt = time.Now()

	// 保存到数据库
	result = h.DB.Save(&workflow)
	if result.Error != nil {
		h.Logger.Error("更新工作流失败", zap.Error(result.Error), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新工作流失败: " + result.Error.Error(),
		})
		return
	}

	// 返回更新后的工作流
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "工作流更新成功",
		"data":    workflow,
	})
}

// DeleteWorkflow 删除工作流
func (h *GinWorkflowHandler) DeleteWorkflow(c *gin.Context) {
	// 获取工作流ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "工作流ID不能为空",
		})
		return
	}

	// 查询工作流是否存在
	var workflow models.Workflow
	result := h.DB.First(&workflow, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "工作流不存在",
			})
		} else {
			h.Logger.Error("查询工作流失败", zap.Error(result.Error), zap.String("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "查询工作流失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 删除工作流
	result = h.DB.Delete(&workflow)
	if result.Error != nil {
		h.Logger.Error("删除工作流失败", zap.Error(result.Error), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "删除工作流失败: " + result.Error.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "工作流删除成功",
	})
}
