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

// CreateRechargeOrderRequest 创建充值订单请求结构
type CreateRechargeOrderRequest struct {
	PackageID     int                  `json:"packageId" binding:"required"`
	PaymentMethod models.PaymentMethod `json:"paymentMethod" binding:"required"`
}

// CreateRechargeOrder 创建充值订单
func (h *GinRechargeHandler) CreateRechargeOrder(c *gin.Context) {
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
	var req CreateRechargeOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求数据格式错误: " + err.Error(),
		})
		return
	}

	// 获取充值套餐信息
	packages := models.GetRechargePackages()
	var selectedPackage *models.RechargePackage
	for _, pkg := range packages {
		if pkg.ID == req.PackageID {
			selectedPackage = pkg
			break
		}
	}

	if selectedPackage == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的充值套餐",
		})
		return
	}

	// 验证支付方式
	validPaymentMethods := []models.PaymentMethod{
		models.PaymentMethodAlipay,
		models.PaymentMethodWechat,
		models.PaymentMethodUnion,
		models.PaymentMethodPaypal,
	}
	validPayment := false
	for _, method := range validPaymentMethods {
		if req.PaymentMethod == method {
			validPayment = true
			break
		}
	}

	if !validPayment {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "不支持的支付方式",
		})
		return
	}

	// 创建充值订单
	order := models.RechargeOrder{
		UserID:        userID.(int64),
		OrderNo:       generateOrderNo(),
		Amount:        selectedPackage.Amount,
		Points:        selectedPackage.TotalPoints,
		PaymentMethod: req.PaymentMethod,
		Status:        models.OrderStatusPending,
	}

	result := h.DB.Create(&order)
	if result.Error != nil {
		h.Logger.Error("创建充值订单失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "创建充值订单失败",
		})
		return
	}

	// 返回订单信息（包含支付相关信息）
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"order": order,
			"payment": gin.H{
				"orderNo":       order.OrderNo,
				"amount":        order.Amount,
				"paymentMethod": order.PaymentMethod,
				"paymentUrl":    generatePaymentUrl(order.OrderNo, order.PaymentMethod),
			},
		},
	})
}

// GetRechargeOrders 获取充值订单列表
func (h *GinRechargeHandler) GetRechargeOrders(c *gin.Context) {
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

	// 获取状态筛选
	statusFilter := c.Query("status")

	// 构建查询
	query := h.DB.Model(&models.RechargeOrder{}).Where("user_id = ?", userID.(int64))

	// 应用状态筛选
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var orders []models.RechargeOrder
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders)
	if result.Error != nil {
		h.Logger.Error("获取充值订单列表失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取充值订单列表失败",
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

// GetRechargeOrder 获取单个充值订单详情
func (h *GinRechargeHandler) GetRechargeOrder(c *gin.Context) {
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

	// 查询订单 - 验证所有权
	var order models.RechargeOrder
	result := h.DB.Where("id = ? AND user_id = ?", id, userID.(int64)).First(&order)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "充值订单不存在",
			})
		} else {
			h.Logger.Error("获取充值订单失败", zap.Error(result.Error), zap.Int64("id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取充值订单失败",
			})
		}
		return
	}

	// 返回订单详情
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   order,
	})
}

// PaymentNotifyRequest 支付通知请求结构
type PaymentNotifyRequest struct {
	OrderNo   string `json:"orderNo" binding:"required"`
	PaymentID string `json:"paymentId" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

// PaymentNotify 处理支付通知
func (h *GinRechargeHandler) PaymentNotify(c *gin.Context) {
	// 解析请求
	var req PaymentNotifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求数据格式错误: " + err.Error(),
		})
		return
	}

	// 查询订单
	var order models.RechargeOrder
	result := h.DB.Where("order_no = ?", req.OrderNo).First(&order)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "订单不存在",
			})
		} else {
			h.Logger.Error("查询充值订单失败", zap.Error(result.Error), zap.String("orderNo", req.OrderNo))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "查询充值订单失败",
			})
		}
		return
	}

	// 检查订单状态
	if order.Status != models.OrderStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "订单状态不正确",
		})
		return
	}

	// 处理支付成功
	if req.Status == "success" || req.Status == "paid" {
		err := h.DB.Transaction(func(tx *gorm.DB) error {
			// 更新订单状态
			if err := tx.Model(&order).Updates(map[string]interface{}{
				"status":     models.OrderStatusPaid,
				"payment_id": req.PaymentID,
				"paid_at":    gorm.Expr("NOW()"),
			}).Error; err != nil {
				return err
			}

			// 获取用户信息
			var user models.User
			if err := tx.First(&user, order.UserID).Error; err != nil {
				return err
			}

			// 增加用户点数
			newBalance := user.Points + order.Points
			if err := tx.Model(&user).Update("points", newBalance).Error; err != nil {
				return err
			}

			// 创建点数交易记录
			transaction := models.PointsTransaction{
				UserID:      order.UserID,
				Type:        models.TransactionTypeRecharge,
				Amount:      order.Points,
				Balance:     newBalance,
				Description: "充值",
				RelatedID:   &order.ID,
			}
			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			h.Logger.Error("处理充值成功通知失败", zap.Error(err), zap.String("orderNo", req.OrderNo))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "处理充值失败",
			})
			return
		}

		h.Logger.Info("充值成功", zap.String("orderNo", req.OrderNo), zap.Int64("userId", order.UserID), zap.Int("points", order.Points))
	} else {
		// 处理支付失败或取消
		status := models.OrderStatusCancelled
		if req.Status == "failed" {
			status = models.OrderStatusCancelled
		}

		if err := h.DB.Model(&order).Updates(map[string]interface{}{
			"status":     status,
			"payment_id": req.PaymentID,
		}).Error; err != nil {
			h.Logger.Error("更新订单状态失败", zap.Error(err), zap.String("orderNo", req.OrderNo))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "更新订单状态失败",
			})
			return
		}
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "处理完成",
	})
}

// PaymentCallbackRequest 支付回调请求结构
type PaymentCallbackRequest struct {
	OrderNumber   string `json:"orderNumber" binding:"required"`
	TransactionID string `json:"transactionId" binding:"required"`
	Status        string `json:"status" binding:"required"`
	Amount        int    `json:"amount" binding:"required"`
	Signature     string `json:"signature" binding:"required"`
}

// PaymentCallback 处理支付回调
func (h *GinRechargeHandler) PaymentCallback(c *gin.Context) {
	// 解析请求体
	var req PaymentCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 验证签名（这里应该实现真实的签名验证逻辑）
	// if !verifySignature(req) {
	//     c.JSON(http.StatusBadRequest, gin.H{
	//         "status": "error",
	//         "message": "签名验证失败",
	//     })
	//     return
	// }

	// 查找充值订单
	var order models.RechargeOrder
	result := h.DB.Where("order_number = ?", req.OrderNumber).First(&order)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "订单不存在",
			})
		} else {
			h.Logger.Error("查询充值订单失败", zap.Error(result.Error), zap.String("orderNumber", req.OrderNumber))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "查询订单失败",
			})
		}
		return
	}

	// 检查订单状态
	if order.Status != models.OrderStatusPending {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "订单状态异常",
		})
		return
	}

	// 验证金额
	if req.Amount != order.Amount {
		h.Logger.Error("支付回调金额不匹配",
			zap.Int("expected", order.Amount),
			zap.Int("received", req.Amount),
			zap.String("orderNumber", req.OrderNumber))
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "金额不匹配",
		})
		return
	}

	// 根据支付状态更新订单
	var orderStatus models.OrderStatus
	switch req.Status {
	case "success", "completed":
		orderStatus = models.OrderStatusPaid
	case "failed", "cancelled":
		orderStatus = models.OrderStatusCancelled
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的支付状态",
		})
		return
	}

	// 使用事务处理充值
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		updates := map[string]interface{}{
			"status":     orderStatus,
			"payment_id": req.TransactionID,
			"paid_at":    time.Now(),
		}

		if err := tx.Model(&order).Updates(updates).Error; err != nil {
			return err
		}

		// 如果支付成功，增加用户点数
		if orderStatus == models.OrderStatusPaid {
			// 获取用户信息
			var user models.User
			if err := tx.First(&user, order.UserID).Error; err != nil {
				return err
			}

			// 增加用户点数
			newBalance := user.Points + order.Points
			if err := tx.Model(&user).Update("points", newBalance).Error; err != nil {
				return err
			}

			// 创建点数交易记录
			transaction := models.PointsTransaction{
				UserID:      order.UserID,
				Type:        models.TransactionTypeRecharge,
				Amount:      order.Points,
				Balance:     newBalance,
				Description: "充值获得点数",
				RelatedID:   &order.ID,
			}

			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		h.Logger.Error("处理支付回调失败", zap.Error(err), zap.String("orderNumber", req.OrderNumber))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "处理支付回调失败",
		})
		return
	}

	// 记录成功日志
	h.Logger.Info("支付回调处理成功",
		zap.String("orderNumber", req.OrderNumber),
		zap.String("transactionId", req.TransactionID),
		zap.String("status", req.Status),
		zap.Int("amount", req.Amount))

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "支付回调处理成功",
	})
}

// generateOrderNo 生成订单号（简化实现）
func generateOrderNo() string {
	return "RO" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

// generatePaymentUrl 生成支付URL（简化实现）
func generatePaymentUrl(orderNo string, method models.PaymentMethod) string {
	// 这里应该根据不同的支付方式生成对应的支付URL
	// 实际实现中需要调用相应的支付服务API
	baseUrl := "https://api.payment.example.com/"
	return baseUrl + "pay?orderNo=" + orderNo + "&method=" + string(method)
}
