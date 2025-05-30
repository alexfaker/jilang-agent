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

// GinRechargeHandler 处理充值相关的请求
type GinRechargeHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewGinRechargeHandler 创建一个新的GinRechargeHandler实例
func NewGinRechargeHandler(db *gorm.DB, logger *zap.Logger) *GinRechargeHandler {
	return &GinRechargeHandler{
		DB:     db,
		Logger: logger,
	}
}

// GetRechargePackages 获取充值套餐列表
func (h *GinRechargeHandler) GetRechargePackages(c *gin.Context) {
	packages := models.GetRechargePackages()
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   packages,
	})
}

// CreateRechargeRequest 创建充值订单请求体
type CreateRechargeRequest struct {
	Amount        int    `json:"amount" binding:"required,min=100"` // 充值金额（分）
	Points        int    `json:"points" binding:"required,min=100"` // 获得点数
	PaymentMethod string `json:"paymentMethod" binding:"required"`  // 支付方式
	PackageID     *int   `json:"packageId,omitempty"`               // 套餐ID（可选）
}

// CreateRecharge 创建充值订单
func (h *GinRechargeHandler) CreateRecharge(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	// 解析请求体
	var req CreateRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数错误：" + err.Error(),
		})
		return
	}

	uid := userID.(string)

	// 验证支付方式
	validPaymentMethods := map[string]models.PaymentMethod{
		"alipay": models.PaymentMethodAlipay,
		"wechat": models.PaymentMethodWechat,
		"credit": models.PaymentMethodCredit,
	}

	paymentMethod, ok := validPaymentMethods[req.PaymentMethod]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "不支持的支付方式",
		})
		return
	}

	// 创建充值订单
	order := &models.RechargeOrder{
		UserID:        uid,
		Amount:        req.Amount,
		Points:        req.Points,
		PaymentMethod: paymentMethod,
		Status:        models.OrderStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 生成订单号
	order.OrderNo = order.GenerateOrderNo()

	// 保存订单到数据库
	if err := h.DB.Create(order).Error; err != nil {
		h.Logger.Error("创建充值订单失败", zap.Error(err), zap.String("userId", uid))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "创建充值订单失败",
		})
		return
	}

	h.Logger.Info("充值订单创建成功",
		zap.String("orderNo", order.OrderNo),
		zap.String("userId", uid),
		zap.Int("amount", req.Amount),
		zap.Int("points", req.Points),
	)

	// 返回订单信息
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"orderId":       order.ID,
			"orderNo":       order.OrderNo,
			"amount":        order.Amount,
			"points":        order.Points,
			"paymentMethod": order.PaymentMethod,
			"status":        order.Status,
			"createdAt":     order.CreatedAt,
			// 在实际环境中，这里应该返回支付链接或支付参数
			"paymentUrl": generatePaymentUrl(order),
		},
	})
}

// GetRechargeHistory 获取充值历史
func (h *GinRechargeHandler) GetRechargeHistory(c *gin.Context) {
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

	uid := userID.(string)

	// 构建查询
	query := h.DB.Model(&models.RechargeOrder{}).Where("user_id = ?", uid)

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var orders []models.RechargeOrder
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders)
	if result.Error != nil {
		h.Logger.Error("获取充值历史失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取充值历史失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"orders": orders,
			"pagination": gin.H{
				"total":     total,
				"page":      offset/limit + 1,
				"page_size": limit,
				"pages":     (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetRechargeStatus 获取充值状态
func (h *GinRechargeHandler) GetRechargeStatus(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "未找到用户信息",
		})
		return
	}

	// 获取订单ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的订单ID",
		})
		return
	}

	uid := userID.(string)

	// 查询订单 - 验证所有权
	var order models.RechargeOrder
	result := h.DB.Where("id = ? AND user_id = ?", id, uid).First(&order)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "订单不存在",
			})
		} else {
			h.Logger.Error("获取订单状态失败", zap.Error(result.Error), zap.Int64("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取订单状态失败",
			})
		}
		return
	}

	// 返回订单状态
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   order,
	})
}

// ProcessPaymentCallback 处理支付回调（用于支付网关回调）
func (h *GinRechargeHandler) ProcessPaymentCallback(c *gin.Context) {
	// 获取订单号
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "订单号不能为空",
		})
		return
	}

	// 查询订单
	var order models.RechargeOrder
	result := h.DB.Where("order_no = ?", orderNo).First(&order)
	if result.Error != nil {
		h.Logger.Error("支付回调：订单不存在", zap.String("orderNo", orderNo))
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "订单不存在",
		})
		return
	}

	// 验证支付状态（这里应该调用支付网关API验证）
	// 模拟支付成功
	if order.Status == models.OrderStatusPending {
		// 开始事务
		tx := h.DB.Begin()

		// 更新订单状态
		order.Status = models.OrderStatusCompleted
		order.PaymentID = "MOCK_PAYMENT_" + strconv.FormatInt(time.Now().Unix(), 10)
		order.UpdatedAt = time.Now()

		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			h.Logger.Error("更新订单状态失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "更新订单状态失败",
			})
			return
		}

		// 更新用户积分
		var user models.User
		if err := tx.Where("user_id = ?", order.UserID).First(&user).Error; err != nil {
			tx.Rollback()
			h.Logger.Error("查询用户失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "查询用户失败",
			})
			return
		}

		user.Points += order.Points
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			h.Logger.Error("更新用户积分失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "更新用户积分失败",
			})
			return
		}

		// 创建积分交易记录
		transaction := &models.PointsTransaction{
			UserID:      order.UserID,
			Type:        models.TransactionTypeRecharge,
			Amount:      order.Points,
			Balance:     user.Points, // 交易后的余额
			Description: "充值获得积分",
			RelatedID:   &order.ID,
			CreatedAt:   time.Now(),
		}

		if err := tx.Create(transaction).Error; err != nil {
			tx.Rollback()
			h.Logger.Error("创建积分交易记录失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "创建积分交易记录失败",
			})
			return
		}

		// 提交事务
		tx.Commit()

		h.Logger.Info("充值支付成功",
			zap.String("orderNo", orderNo),
			zap.String("userId", order.UserID),
			zap.Int("points", order.Points),
		)
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   order,
	})
}

// generatePaymentUrl 生成支付链接（模拟）
func generatePaymentUrl(order *models.RechargeOrder) string {
	// 在实际环境中，这里应该调用支付网关API生成支付链接
	return "https://pay.example.com/pay?orderNo=" + order.OrderNo
}
