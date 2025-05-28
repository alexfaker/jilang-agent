package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GinAgentHandler 处理工作流商店相关的请求
type GinAgentHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewGinAgentHandler 创建一个新的GinAgentHandler实例
func NewGinAgentHandler(db *gorm.DB, logger *zap.Logger) *GinAgentHandler {
	return &GinAgentHandler{
		DB:     db,
		Logger: logger,
	}
}

// GetAgents 获取工作流商店中的代理列表
func (h *GinAgentHandler) GetAgents(c *gin.Context) {
	// 解析查询参数
	category := c.Query("category")
	search := c.Query("search")

	// 分页参数
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

	// 查询代理列表（只查询公开的代理）
	var agents []models.Agent
	query := h.DB.Model(&models.Agent{}).Where("is_public = ?", true)

	// 应用类别筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 应用搜索筛选
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 获取分页数据（按购买次数、评分、创建时间排序）
	result := query.Order("purchase_count DESC, rating DESC, created_at DESC").Offset(offset).Limit(limit).Find(&agents)
	if result.Error != nil {
		h.Logger.Error("获取代理列表失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取代理列表失败: " + result.Error.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"agents": agents,
			"pagination": gin.H{
				"total":     total,
				"page":      offset/limit + 1,
				"page_size": limit,
				"pages":     (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetAgent 获取单个代理详情
func (h *GinAgentHandler) GetAgent(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的代理ID",
		})
		return
	}

	// 获取代理（只能获取公开的代理）
	var agent models.Agent
	result := h.DB.Where("is_public = ?", true).First(&agent, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "代理不存在",
			})
		} else {
			h.Logger.Error("获取代理失败", zap.Error(result.Error), zap.Int64("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取代理失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   agent,
	})
}

// GetAgentCategories 获取代理分类列表
func (h *GinAgentHandler) GetAgentCategories(c *gin.Context) {
	// 查询公开代理的所有不同分类
	var categories []string
	result := h.DB.Model(&models.Agent{}).Where("is_public = ? AND category != ''", true).Distinct().Pluck("category", &categories)
	if result.Error != nil {
		h.Logger.Error("获取代理分类失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取代理分类失败: " + result.Error.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   categories,
	})
}

// AgentCreateRequest 创建代理请求结构
type AgentCreateRequest struct {
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	Type        string          `json:"type" binding:"required"`
	Category    string          `json:"category" binding:"required"`
	Icon        string          `json:"icon"`
	Definition  json.RawMessage `json:"definition" binding:"required"`
	Price       int             `json:"price" binding:"required,min=0"`
	IsPublic    bool            `json:"is_public"`
}

// CreateAgent 创建代理（管理员功能）
func (h *GinAgentHandler) CreateAgent(c *gin.Context) {
	// 检查管理员权限
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "权限不足，只有管理员可以创建代理",
		})
		return
	}

	// 解析请求体
	var req AgentCreateRequest
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
			"message": "代理定义必须是有效的JSON格式",
		})
		return
	}

	// 创建代理
	agent := models.Agent{
		Name:          req.Name,
		Description:   req.Description,
		Type:          req.Type,
		Category:      req.Category,
		Icon:          req.Icon,
		Definition:    req.Definition,
		Price:         req.Price,
		PurchaseCount: 0,
		Rating:        0.0,
		IsPublic:      req.IsPublic,
	}

	result := h.DB.Create(&agent)
	if result.Error != nil {
		h.Logger.Error("创建代理失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "创建代理失败: " + result.Error.Error(),
		})
		return
	}

	// 返回创建的代理
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   agent,
	})
}

// AgentUpdateRequest 更新代理请求结构
type AgentUpdateRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Category    string          `json:"category"`
	Icon        string          `json:"icon"`
	Definition  json.RawMessage `json:"definition"`
	Price       int             `json:"price"`
	IsPublic    *bool           `json:"is_public"`
}

// UpdateAgent 更新代理（管理员功能）
func (h *GinAgentHandler) UpdateAgent(c *gin.Context) {
	// 检查管理员权限
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "权限不足，只有管理员可以更新代理",
		})
		return
	}

	// 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的代理ID",
		})
		return
	}

	// 检查代理是否存在
	var agent models.Agent
	result := h.DB.First(&agent, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "代理不存在",
			})
		} else {
			h.Logger.Error("获取代理失败", zap.Error(result.Error), zap.Int64("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取代理失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 解析请求体
	var req AgentUpdateRequest
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
				"message": "代理定义必须是有效的JSON格式",
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
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if len(req.Definition) > 0 {
		updates["definition"] = req.Definition
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	// 更新代理
	result = h.DB.Model(&agent).Updates(updates)
	if result.Error != nil {
		h.Logger.Error("更新代理失败", zap.Error(result.Error), zap.Int64("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新代理失败: " + result.Error.Error(),
		})
		return
	}

	// 重新获取更新后的代理
	h.DB.First(&agent, id)

	// 返回更新后的代理
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   agent,
	})
}

// DeleteAgent 删除代理（管理员功能）
func (h *GinAgentHandler) DeleteAgent(c *gin.Context) {
	// 检查管理员权限
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "权限不足，只有管理员可以删除代理",
		})
		return
	}

	// 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的代理ID",
		})
		return
	}

	// 检查代理是否存在
	var agent models.Agent
	result := h.DB.First(&agent, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "代理不存在",
			})
		} else {
			h.Logger.Error("获取代理失败", zap.Error(result.Error), zap.Int64("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取代理失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 检查是否有用户购买了此代理
	var workflowCount int64
	h.DB.Model(&models.Workflow{}).Where("agent_id = ?", id).Count(&workflowCount)
	if workflowCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "无法删除代理，已有用户购买了此代理",
		})
		return
	}

	// 删除代理
	result = h.DB.Delete(&agent)
	if result.Error != nil {
		h.Logger.Error("删除代理失败", zap.Error(result.Error), zap.Int64("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "删除代理失败: " + result.Error.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "代理删除成功",
	})
}
