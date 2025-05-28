package handlers

import (
	"net/http"
	"strconv"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GinPointsHandler 处理点数相关的请求
type GinPointsHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewGinPointsHandler 创建一个新的GinPointsHandler实例
func NewGinPointsHandler(db *gorm.DB, logger *zap.Logger) *GinPointsHandler {
	return &GinPointsHandler{
		DB:     db,
		Logger: logger,
	}
}

// GetPointsBalance 获取用户点数余额
func (h *GinPointsHandler) GetPointsBalance(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	// 查询用户信息
	var user models.User
	result := h.DB.First(&user, userID.(int64))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户信息失败", zap.Error(result.Error), zap.Int64("userId", userID.(int64)))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户信息失败",
			})
		}
		return
	}

	// 返回余额信息
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"points": user.Points,
			"userId": user.ID,
		},
	})
}

// GetPointsTransactions 获取点数交易历史
func (h *GinPointsHandler) GetPointsTransactions(c *gin.Context) {
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

	// 获取类型筛选
	typeFilter := c.Query("type")

	// 构建查询
	query := h.DB.Model(&models.PointsTransaction{}).Where("user_id = ?", userID.(int64))

	// 应用类型筛选
	if typeFilter != "" {
		query = query.Where("type = ?", typeFilter)
	}

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var transactions []models.PointsTransaction
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions)
	if result.Error != nil {
		h.Logger.Error("获取交易历史失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取交易历史失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"transactions": transactions,
			"pagination": gin.H{
				"total":     total,
				"page":      offset/limit + 1,
				"page_size": limit,
				"pages":     (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetPointsTransaction 获取单个交易详情
func (h *GinPointsHandler) GetPointsTransaction(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	// 获取交易ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的交易ID",
		})
		return
	}

	// 查询交易 - 验证所有权
	var transaction models.PointsTransaction
	result := h.DB.Where("id = ? AND user_id = ?", id, userID.(int64)).First(&transaction)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "交易记录不存在",
			})
		} else {
			h.Logger.Error("获取交易记录失败", zap.Error(result.Error), zap.Int64("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取交易记录失败",
			})
		}
		return
	}

	// 返回交易详情
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transaction,
	})
}

// GetPointsStatistics 获取点数统计信息
func (h *GinPointsHandler) GetPointsStatistics(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	uid := userID.(int64)

	// 查询用户当前余额
	var user models.User
	if err := h.DB.First(&user, uid).Error; err != nil {
		h.Logger.Error("获取用户信息失败", zap.Error(err), zap.Int64("userId", uid))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取用户信息失败",
		})
		return
	}

	// 统计总充值
	var totalRecharge int64
	h.DB.Model(&models.PointsTransaction{}).
		Where("user_id = ? AND type = ?", uid, models.TransactionTypeRecharge).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalRecharge)

	// 统计总消费
	var totalSpent int64
	h.DB.Model(&models.PointsTransaction{}).
		Where("user_id = ? AND type IN ?", uid, []models.TransactionType{
			models.TransactionTypePurchase,
			models.TransactionTypeExecution,
		}).
		Select("COALESCE(SUM(ABS(amount)), 0)").
		Scan(&totalSpent)

	// 统计购买的工作流数量
	var workflowCount int64
	h.DB.Model(&models.Workflow{}).
		Where("user_id = ? AND agent_id IS NOT NULL", uid).
		Count(&workflowCount)

	// 统计最近30天的交易
	var recentTransactionCount int64
	h.DB.Model(&models.PointsTransaction{}).
		Where("user_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)", uid).
		Count(&recentTransactionCount)

	// 返回统计信息
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"currentBalance":         user.Points,
			"totalRecharge":          totalRecharge,
			"totalSpent":             totalSpent,
			"purchasedWorkflowCount": workflowCount,
			"recentTransactionCount": recentTransactionCount,
		},
	})
}
