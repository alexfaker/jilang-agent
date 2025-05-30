package models

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TransactionType 交易类型
type TransactionType string

const (
	TransactionTypeRecharge  TransactionType = "recharge"  // 充值
	TransactionTypePurchase  TransactionType = "purchase"  // 购买
	TransactionTypeExecution TransactionType = "execution" // 执行消费
	TransactionTypeRefund    TransactionType = "refund"    // 退款
)

// PointsTransaction 点数交易记录模型
type PointsTransaction struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      string          `json:"userID" gorm:"column:user_id;index;not null"`
	Type        TransactionType `json:"type" gorm:"type:varchar(20);not null"`
	Amount      int             `json:"amount" gorm:"not null"`                   // 正数为增加，负数为减少
	Balance     int             `json:"balance" gorm:"not null"`                  // 交易后余额
	Description string          `json:"description" gorm:"type:text"`             // 交易描述
	RelatedID   *int64          `json:"relatedId" gorm:"column:related_id;index"` // 关联ID（工作流ID、订单ID等）
	CreatedAt   time.Time       `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

// TableName 指定表名
func (PointsTransaction) TableName() string {
	return "points_transactions"
}

// PointsTransactionCreateInput 创建点数交易输入
type PointsTransactionCreateInput struct {
	UserID      string          `json:"userID" validate:"required"`
	Type        TransactionType `json:"type" validate:"required"`
	Amount      int             `json:"amount" validate:"required"`
	Description string          `json:"description"`
	RelatedID   *int64          `json:"relatedId"`
}

// CreatePointsTransaction 创建点数交易记录
func CreatePointsTransaction(db *sql.DB, input PointsTransactionCreateInput) (*PointsTransaction, error) {
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("开始事务失败: %w", err)
	}
	defer tx.Rollback()

	// 获取用户当前余额
	var currentBalance int
	err = tx.QueryRow("SELECT points FROM users WHERE id = ?", input.UserID).Scan(&currentBalance)
	if err != nil {
		return nil, fmt.Errorf("获取用户余额失败: %w", err)
	}

	// 计算新余额
	newBalance := currentBalance + input.Amount
	if newBalance < 0 {
		return nil, fmt.Errorf("余额不足，当前余额: %d，尝试扣除: %d", currentBalance, -input.Amount)
	}

	// 更新用户余额
	_, err = tx.Exec("UPDATE users SET points = ? WHERE id = ?", newBalance, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("更新用户余额失败: %w", err)
	}

	// 创建交易记录
	query := `
		INSERT INTO points_transactions (user_id, type, amount, balance, description, related_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
		RETURNING id, created_at
	`

	var transaction PointsTransaction
	transaction.UserID = input.UserID
	transaction.Type = input.Type
	transaction.Amount = input.Amount
	transaction.Balance = newBalance
	transaction.Description = input.Description
	transaction.RelatedID = input.RelatedID

	err = tx.QueryRow(
		query,
		transaction.UserID,
		transaction.Type,
		transaction.Amount,
		transaction.Balance,
		transaction.Description,
		transaction.RelatedID,
	).Scan(&transaction.ID, &transaction.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("创建交易记录失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	return &transaction, nil
}

// ListPointsTransactions 获取用户的点数交易历史
func ListPointsTransactions(db *sql.DB, userID int64, transactionType *TransactionType, limit, offset int) ([]*PointsTransaction, error) {
	query := `
		SELECT id, user_id, type, amount, balance, description, related_id, created_at
		FROM points_transactions
		WHERE user_id = ?
	`
	args := []interface{}{userID}

	// 添加类型筛选
	if transactionType != nil {
		query += " AND type = ?"
		args = append(args, string(*transactionType))
	}

	// 添加排序和分页
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询交易历史失败: %w", err)
	}
	defer rows.Close()

	var transactions []*PointsTransaction
	for rows.Next() {
		var transaction PointsTransaction
		var relatedID sql.NullInt64

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Balance,
			&transaction.Description,
			&relatedID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描交易记录失败: %w", err)
		}

		if relatedID.Valid {
			transaction.RelatedID = &relatedID.Int64
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// GetUserPointsBalance 获取用户点数余额
func GetUserPointsBalance(db *sql.DB, userID int64) (int, error) {
	var balance int
	err := db.QueryRow("SELECT points FROM users WHERE id = ?", userID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("用户不存在")
		}
		return 0, fmt.Errorf("获取用户余额失败: %w", err)
	}
	return balance, nil
}

// CreatePointsTransactionGorm 使用GORM创建点数交易记录
func CreatePointsTransactionGorm(db *gorm.DB, input PointsTransactionCreateInput) (*PointsTransaction, error) {
	var transaction PointsTransaction

	err := db.Transaction(func(tx *gorm.DB) error {
		// 获取用户当前余额
		var user User
		if err := tx.First(&user, input.UserID).Error; err != nil {
			return fmt.Errorf("获取用户信息失败: %w", err)
		}

		// 计算新余额
		newBalance := user.Points + input.Amount
		if newBalance < 0 {
			return fmt.Errorf("余额不足，当前余额: %d，尝试扣除: %d", user.Points, -input.Amount)
		}

		// 更新用户余额
		if err := tx.Model(&user).Update("points", newBalance).Error; err != nil {
			return fmt.Errorf("更新用户余额失败: %w", err)
		}

		// 创建交易记录
		transaction = PointsTransaction{
			UserID:      input.UserID,
			Type:        input.Type,
			Amount:      input.Amount,
			Balance:     newBalance,
			Description: input.Description,
			RelatedID:   input.RelatedID,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("创建交易记录失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// ListPointsTransactionsGorm 使用GORM获取用户的点数交易历史
func ListPointsTransactionsGorm(db *gorm.DB, userID int64, transactionType *TransactionType, limit, offset int) ([]*PointsTransaction, error) {
	var transactions []*PointsTransaction
	query := db.Where("user_id = ?", userID)

	// 添加类型筛选
	if transactionType != nil {
		query = query.Where("type = ?", *transactionType)
	}

	// 执行查询
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	if err != nil {
		return nil, fmt.Errorf("查询交易历史失败: %w", err)
	}

	return transactions, nil
}
