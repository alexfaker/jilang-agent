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

// GinAgentHandler 处理代理相关的请求
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

// GetAgents 获取代理列表
func (h *GinAgentHandler) GetAgents(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}

	// 解析查询参数
	category := c.Query("category")

	// 是否公开筛选（可选）
	var isPublic *bool
	if publicStr := c.Query("is_public"); publicStr != "" {
		if publicStr == "true" {
			trueValue := true
			isPublic = &trueValue
		} else if publicStr == "false" {
			falseValue := false
			isPublic = &falseValue
		}
	}

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

	// 查询代理列表
	var agents []models.Agent
	query := h.DB.Model(&models.Agent{})

	// 应用筛选条件
	uid := userID.(int64)
	if isPublic != nil {
		if *isPublic {
			// 查询公开的代理
			query = query.Where("is_public = ?", true)
		} else {
			// 查询用户自己的非公开代理
			query = query.Where("user_id = ?", uid).Where("is_public = ?", false)
		}
	} else {
		// 查询公开的代理和用户自己的代理
		query = query.Where("is_public = ? OR user_id = ?", true, uid)
	}

	// 应用类别筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 获取分页数据
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&agents)
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

// GetAgent 获取单个代理
func (h *GinAgentHandler) GetAgent(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
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

	// 获取代理
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

	// 验证用户权限（只能查看公开的代理或自己的代理）
	uid := userID.(int64)
	if !agent.IsPublic && (agent.UserID == nil || *agent.UserID != uid) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "无权访问此代理",
		})
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
	// 查询所有不同的代理分类
	var categories []string
	result := h.DB.Model(&models.Agent{}).Distinct().Pluck("category", &categories)
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
	IsPublic    bool            `json:"is_public"`
}

// CreateAgent 创建代理
func (h *GinAgentHandler) CreateAgent(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}

	// 解析请求体
	var req AgentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 验证Definition是有效的JSON
	var js json.RawMessage
	if err := json.Unmarshal(req.Definition, &js); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "代理定义不是有效的JSON: " + err.Error(),
		})
		return
	}

	// 创建代理
	uid := userID.(int64)
	agent := models.Agent{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Category:    req.Category,
		Icon:        req.Icon,
		Definition:  req.Definition,
		IsPublic:    req.IsPublic,
		UserID:      &uid,
	}

	// 保存到数据库
	result := h.DB.Create(&agent)
	if result.Error != nil {
		h.Logger.Error("创建代理失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "创建代理失败: " + result.Error.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "代理创建成功",
		"data":    agent,
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
	IsPublic    *bool           `json:"is_public"`
}

// UpdateAgent 更新代理
func (h *GinAgentHandler) UpdateAgent(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
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

	// 获取代理
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

	// 验证用户权限（只能更新自己的代理）
	uid := userID.(int64)
	if agent.UserID == nil || *agent.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "无权更新此代理",
		})
		return
	}

	// 解析请求体
	var req AgentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 验证Definition是有效的JSON（如果提供）
	if len(req.Definition) > 0 {
		var js json.RawMessage
		if err := json.Unmarshal(req.Definition, &js); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "代理定义不是有效的JSON: " + err.Error(),
			})
			return
		}
	}

	// 更新代理
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
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

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

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "代理更新成功",
		"data":    agent,
	})
}

// DeleteAgent 删除代理
func (h *GinAgentHandler) DeleteAgent(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
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

	// 获取代理
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

	// 验证用户权限（只能删除自己的代理）
	uid := userID.(int64)
	if agent.UserID == nil || *agent.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "无权删除此代理",
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

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "代理删除成功",
	})
}
