package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Agent 代理模型 - 工作流商店中的模板
type Agent struct {
	ID            int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string          `json:"name" gorm:"type:varchar(100);not null"`
	Description   string          `json:"description" gorm:"type:text"`
	Type          string          `json:"type" gorm:"type:varchar(50);not null"`
	Category      string          `json:"category" gorm:"type:varchar(50);index"`
	Icon          string          `json:"icon" gorm:"type:varchar(255)"`
	CoverImage    string          `json:"coverImage" gorm:"column:cover_image;type:varchar(500)"` // 封面图URL
	Definition    json.RawMessage `json:"definition" gorm:"type:json"`                            // 代理定义JSON
	Price         int             `json:"price" gorm:"not null;default:0"`                        // 价格（点数）
	PurchaseCount int             `json:"purchaseCount" gorm:"column:purchase_count;default:0"`   // 购买次数
	Rating        float64         `json:"rating" gorm:"default:0.0"`                              // 评分
	IsPublic      bool            `json:"isPublic" gorm:"column:is_public;default:false"`
	CreatedAt     time.Time       `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time       `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (Agent) TableName() string {
	return "agents"
}

// AgentCreateInput 创建代理输入
type AgentCreateInput struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Type        string          `json:"type" validate:"required"`
	Category    string          `json:"category" validate:"required"`
	Icon        string          `json:"icon"`
	CoverImage  string          `json:"coverImage"`
	Definition  json.RawMessage `json:"definition" validate:"required"`
	Price       int             `json:"price" validate:"required,min=0"`
	IsPublic    bool            `json:"isPublic"`
}

// AgentUpdateInput 更新代理输入
type AgentUpdateInput struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Category    string          `json:"category"`
	Icon        string          `json:"icon"`
	CoverImage  string          `json:"coverImage"`
	Definition  json.RawMessage `json:"definition"`
	Price       int             `json:"price"`
	IsPublic    bool            `json:"isPublic"`
}

// CreateAgent 创建新代理
func CreateAgent(db *gorm.DB, input AgentCreateInput) (*Agent, error) {
	agent := Agent{
		Name:          input.Name,
		Description:   input.Description,
		Type:          input.Type,
		Category:      input.Category,
		Icon:          input.Icon,
		CoverImage:    input.CoverImage,
		Definition:    input.Definition,
		Price:         input.Price,
		PurchaseCount: 0,
		Rating:        0.0,
		IsPublic:      input.IsPublic,
	}

	if err := db.Create(&agent).Error; err != nil {
		return nil, err
	}

	return &agent, nil
}

// GetAgent 获取代理
func GetAgent(db *gorm.DB, id int64) (*Agent, error) {
	var agent Agent
	if err := db.First(&agent, id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// ListAgents 获取代理列表（工作流商店）
func ListAgents(db *gorm.DB, category string, search string, limit, offset int) ([]*Agent, error) {
	var agents []*Agent
	query := db.Where("is_public = ?", true)

	// 添加筛选条件
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 添加排序和分页
	err := query.Order("purchase_count DESC, rating DESC, created_at DESC").
		Offset(offset).Limit(limit).Find(&agents).Error
	if err != nil {
		return nil, err
	}

	return agents, nil
}

// Update 更新代理
func (a *Agent) Update(db *gorm.DB, input AgentUpdateInput) error {
	updates := map[string]interface{}{}

	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Type != "" {
		updates["type"] = input.Type
	}
	if input.Category != "" {
		updates["category"] = input.Category
	}
	if input.Icon != "" {
		updates["icon"] = input.Icon
	}
	if input.CoverImage != "" {
		updates["cover_image"] = input.CoverImage
	}
	if len(input.Definition) > 0 {
		updates["definition"] = input.Definition
	}
	if input.Price > 0 {
		updates["price"] = input.Price
	}
	updates["is_public"] = input.IsPublic

	return db.Model(a).Updates(updates).Error
}

// Delete 删除代理
func (a *Agent) Delete(db *gorm.DB) error {
	return db.Delete(a).Error
}

// IncrementPurchaseCount 增加购买次数
func (a *Agent) IncrementPurchaseCount(db *gorm.DB) error {
	return db.Model(a).Update("purchase_count", gorm.Expr("purchase_count + ?", 1)).Error
}

// GetAgentCategories 获取代理分类列表
func GetAgentCategories(db *gorm.DB) ([]string, error) {
	var categories []string
	err := db.Model(&Agent{}).Where("is_public = ? AND category != ''", true).
		Distinct().Pluck("category", &categories).Error
	return categories, err
}
