package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GinPurchaseHandler 处理购买相关的请求
type GinPurchaseHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewGinPurchaseHandler 创建一个新的GinPurchaseHandler实例
func NewGinPurchaseHandler(db *gorm.DB, logger *zap.Logger) *GinPurchaseHandler {
	return &GinPurchaseHandler{
		DB:     db,
		Logger: logger,
	}
}

// PurchaseAgentRequest 购买代理请求结构
type PurchaseAgentRequest struct {
	AgentID int64 `json:"agentId" binding:"required"`
}

// PurchaseError 购买错误
type PurchaseError struct {
	Message string
}

func (e *PurchaseError) Error() string {
	return e.Message
}

// PurchaseAgent 购买代理
func (h *GinPurchaseHandler) PurchaseAgent(c *gin.Context) {
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
	var req PurchaseAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求数据格式错误: " + err.Error(),
		})
		return
	}

	uid := userID.(int64)

	// 使用事务处理购买流程
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 获取代理信息
		var agent models.Agent
		if err := tx.Where("id = ? AND is_public = ?", req.AgentID, true).First(&agent).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &PurchaseError{Message: "代理不存在或不可购买"}
			}
			return err
		}

		// 检查用户是否已经购买过此代理
		var existingWorkflow models.Workflow
		err := tx.Where("user_id = ? AND agent_id = ?", uid, req.AgentID).First(&existingWorkflow).Error
		if err == nil {
			return &PurchaseError{Message: "您已经购买过此代理"}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// 获取用户信息
		var user models.User
		if err := tx.First(&user, uid).Error; err != nil {
			return err
		}

		// 检查用户余额
		if user.Points < agent.Price {
			return &PurchaseError{Message: "余额不足，请先充值"}
		}

		// 扣除用户点数
		if err := tx.Model(&user).Update("points", user.Points-agent.Price).Error; err != nil {
			return err
		}

		// 创建点数交易记录
		transaction := models.PointsTransaction{
			UserID:      uid,
			Type:        models.TransactionTypePurchase,
			Amount:      -agent.Price,
			Balance:     user.Points - agent.Price,
			Description: "购买工作流: " + agent.Name,
			RelatedID:   &agent.ID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// 创建工作流实例
		now := time.Now()
		workflow := models.Workflow{
			Name:        agent.Name,
			Description: agent.Description,
			UserID:      uid,
			AgentID:     &agent.ID,
			Status:      models.WorkflowStatusActive,
			Definition:  agent.Definition,
			PurchasedAt: &now,
			RunCount:    0,
		}
		if err := tx.Create(&workflow).Error; err != nil {
			return err
		}

		// 增加代理购买次数
		if err := tx.Model(&agent).Update("purchase_count", gorm.Expr("purchase_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if purchaseErr, ok := err.(*PurchaseError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": purchaseErr.Message,
			})
		} else {
			h.Logger.Error("购买代理失败", zap.Error(err), zap.Int64("agentId", req.AgentID), zap.Int64("userId", uid))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "购买代理失败",
			})
		}
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "购买成功",
	})
}

// GetPurchaseHistory 获取购买历史
func (h *GinPurchaseHandler) GetPurchaseHistory(c *gin.Context) {
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
				limit = 100
			}
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 查询用户购买的工作流（只查询从代理购买的）
	var workflows []models.Workflow
	query := h.DB.Where("user_id = ? AND agent_id IS NOT NULL", userID.(int64))

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 获取分页数据
	result := query.Order("purchased_at DESC").Offset(offset).Limit(limit).Find(&workflows)
	if result.Error != nil {
		h.Logger.Error("获取购买历史失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取购买历史失败",
		})
		return
	}

	// 为每个工作流获取关联的代理信息
	type PurchaseHistoryItem struct {
		models.Workflow
		Agent *models.Agent `json:"agent,omitempty"`
	}

	var purchaseHistory []PurchaseHistoryItem
	for _, workflow := range workflows {
		item := PurchaseHistoryItem{Workflow: workflow}

		// 如果有代理ID，获取代理信息
		if workflow.AgentID != nil {
			var agent models.Agent
			if err := h.DB.First(&agent, *workflow.AgentID).Error; err == nil {
				item.Agent = &agent
			}
		}

		purchaseHistory = append(purchaseHistory, item)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"purchases": purchaseHistory,
			"pagination": gin.H{
				"total":     total,
				"page":      offset/limit + 1,
				"page_size": limit,
				"pages":     (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}
