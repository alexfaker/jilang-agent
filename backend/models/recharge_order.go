package models

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusPaid      OrderStatus = "paid"      // 已支付
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
)

// PaymentMethod 支付方式
type PaymentMethod string

const (
	PaymentMethodAlipay PaymentMethod = "alipay" // 支付宝
	PaymentMethodWechat PaymentMethod = "wechat" // 微信支付
	PaymentMethodUnion  PaymentMethod = "union"  // 银联
	PaymentMethodPaypal PaymentMethod = "paypal" // PayPal
)

// RechargeOrder 充值订单模型
type RechargeOrder struct {
	ID            int64         `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int64         `json:"userId" gorm:"column:user_id;index;not null"`
	OrderNo       string        `json:"orderNo" gorm:"column:order_no;uniqueIndex;not null"` // 订单号
	Amount        int           `json:"amount" gorm:"not null"`                              // 充值金额（分）
	Points        int           `json:"points" gorm:"not null"`                              // 获得点数
	PaymentMethod PaymentMethod `json:"paymentMethod" gorm:"column:payment_method;type:varchar(20);not null"`
	Status        OrderStatus   `json:"status" gorm:"type:varchar(20);default:'pending';not null"`
	PaymentID     string        `json:"paymentId" gorm:"column:payment_id;index"` // 第三方支付ID
	PaidAt        *time.Time    `json:"paidAt" gorm:"column:paid_at"`             // 支付时间
	CreatedAt     time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (RechargeOrder) TableName() string {
	return "recharge_orders"
}

// RechargePackage 充值套餐
type RechargePackage struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`      // 价格（分）
	Points      int    `json:"points"`      // 基础点数
	BonusPoints int    `json:"bonusPoints"` // 赠送点数
	TotalPoints int    `json:"totalPoints"` // 总点数
	Popular     bool   `json:"popular"`     // 是否热门
}

// RechargeOrderCreateInput 创建充值订单输入
type RechargeOrderCreateInput struct {
	UserID        int64         `json:"userId" validate:"required"`
	Amount        int           `json:"amount" validate:"required,min=1"`
	Points        int           `json:"points" validate:"required,min=1"`
	PaymentMethod PaymentMethod `json:"paymentMethod" validate:"required"`
}

// CreateRechargeOrder 创建充值订单
func CreateRechargeOrder(db *sql.DB, input RechargeOrderCreateInput) (*RechargeOrder, error) {
	// 生成订单号
	orderNo := generateOrderNo()

	query := `
		INSERT INTO recharge_orders (user_id, order_no, amount, points, payment_method, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 'pending', NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var order RechargeOrder
	order.UserID = input.UserID
	order.OrderNo = orderNo
	order.Amount = input.Amount
	order.Points = input.Points
	order.PaymentMethod = input.PaymentMethod
	order.Status = OrderStatusPending

	err := db.QueryRow(
		query,
		order.UserID,
		order.OrderNo,
		order.Amount,
		order.Points,
		order.PaymentMethod,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("创建充值订单失败: %w", err)
	}

	return &order, nil
}

// GetRechargeOrder 获取充值订单
func GetRechargeOrder(db *sql.DB, id int64) (*RechargeOrder, error) {
	query := `
		SELECT id, user_id, order_no, amount, points, payment_method, status, payment_id, paid_at, created_at, updated_at
		FROM recharge_orders
		WHERE id = ?
	`

	var order RechargeOrder
	var paymentID sql.NullString
	var paidAt sql.NullTime

	err := db.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNo,
		&order.Amount,
		&order.Points,
		&order.PaymentMethod,
		&order.Status,
		&paymentID,
		&paidAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	if paymentID.Valid {
		order.PaymentID = paymentID.String
	}

	if paidAt.Valid {
		order.PaidAt = &paidAt.Time
	}

	return &order, nil
}

// GetRechargeOrderByOrderNo 根据订单号获取充值订单
func GetRechargeOrderByOrderNo(db *sql.DB, orderNo string) (*RechargeOrder, error) {
	query := `
		SELECT id, user_id, order_no, amount, points, payment_method, status, payment_id, paid_at, created_at, updated_at
		FROM recharge_orders
		WHERE order_no = ?
	`

	var order RechargeOrder
	var paymentID sql.NullString
	var paidAt sql.NullTime

	err := db.QueryRow(query, orderNo).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNo,
		&order.Amount,
		&order.Points,
		&order.PaymentMethod,
		&order.Status,
		&paymentID,
		&paidAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	if paymentID.Valid {
		order.PaymentID = paymentID.String
	}

	if paidAt.Valid {
		order.PaidAt = &paidAt.Time
	}

	return &order, nil
}

// UpdateOrderStatus 更新订单状态
func (o *RechargeOrder) UpdateStatus(db *sql.DB, status OrderStatus, paymentID string) error {
	query := `
		UPDATE recharge_orders
		SET status = ?, payment_id = ?, paid_at = ?, updated_at = NOW()
		WHERE id = ?
	`

	var paidAt interface{}
	if status == OrderStatusPaid {
		now := time.Now()
		o.PaidAt = &now
		paidAt = now
	} else {
		paidAt = nil
	}

	_, err := db.Exec(query, status, paymentID, paidAt, o.ID)
	if err != nil {
		return fmt.Errorf("更新订单状态失败: %w", err)
	}

	o.Status = status
	o.PaymentID = paymentID
	o.UpdatedAt = time.Now()

	return nil
}

// ListRechargeOrders 获取用户的充值订单列表
func ListRechargeOrders(db *sql.DB, userID int64, status *OrderStatus, limit, offset int) ([]*RechargeOrder, error) {
	query := `
		SELECT id, user_id, order_no, amount, points, payment_method, status, payment_id, paid_at, created_at, updated_at
		FROM recharge_orders
		WHERE user_id = ?
	`
	args := []interface{}{userID}

	// 添加状态筛选
	if status != nil {
		query += " AND status = ?"
		args = append(args, string(*status))
	}

	// 添加排序和分页
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询充值订单失败: %w", err)
	}
	defer rows.Close()

	var orders []*RechargeOrder
	for rows.Next() {
		var order RechargeOrder
		var paymentID sql.NullString
		var paidAt sql.NullTime

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNo,
			&order.Amount,
			&order.Points,
			&order.PaymentMethod,
			&order.Status,
			&paymentID,
			&paidAt,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描订单数据失败: %w", err)
		}

		if paymentID.Valid {
			order.PaymentID = paymentID.String
		}

		if paidAt.Valid {
			order.PaidAt = &paidAt.Time
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

// GetRechargePackages 获取充值套餐列表
func GetRechargePackages() []*RechargePackage {
	return []*RechargePackage{
		{
			ID:          1,
			Name:        "入门套餐",
			Amount:      1000, // 10元
			Points:      1000,
			BonusPoints: 0,
			TotalPoints: 1000,
			Popular:     false,
		},
		{
			ID:          2,
			Name:        "标准套餐",
			Amount:      5000, // 50元
			Points:      5000,
			BonusPoints: 500,
			TotalPoints: 5500,
			Popular:     true,
		},
		{
			ID:          3,
			Name:        "高级套餐",
			Amount:      10000, // 100元
			Points:      10000,
			BonusPoints: 2000,
			TotalPoints: 12000,
			Popular:     false,
		},
		{
			ID:          4,
			Name:        "专业套餐",
			Amount:      20000, // 200元
			Points:      20000,
			BonusPoints: 5000,
			TotalPoints: 25000,
			Popular:     false,
		},
	}
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	return fmt.Sprintf("RO%d", time.Now().UnixNano())
}

// CreateRechargeOrderGorm 使用GORM创建充值订单
func CreateRechargeOrderGorm(db *gorm.DB, input RechargeOrderCreateInput) (*RechargeOrder, error) {
	order := &RechargeOrder{
		UserID:        input.UserID,
		OrderNo:       generateOrderNo(),
		Amount:        input.Amount,
		Points:        input.Points,
		PaymentMethod: input.PaymentMethod,
		Status:        OrderStatusPending,
	}

	if err := db.Create(order).Error; err != nil {
		return nil, fmt.Errorf("创建充值订单失败: %w", err)
	}

	return order, nil
}

// UpdateStatusGorm 使用GORM更新订单状态
func (o *RechargeOrder) UpdateStatusGorm(db *gorm.DB, status OrderStatus, paymentID string) error {
	updates := map[string]interface{}{
		"status":     status,
		"payment_id": paymentID,
	}

	if status == OrderStatusPaid {
		now := time.Now()
		o.PaidAt = &now
		updates["paid_at"] = now
	}

	if err := db.Model(o).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新订单状态失败: %w", err)
	}

	o.Status = status
	o.PaymentID = paymentID

	return nil
}
