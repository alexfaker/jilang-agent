package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Agent 代理模型
type Agent struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Category    string          `json:"category"`
	Icon        string          `json:"icon"`
	Definition  json.RawMessage `json:"definition"` // 代理定义JSON
	IsPublic    bool            `json:"isPublic"`
	UserID      *int64          `json:"userId"` // 为空表示系统预置
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	UsageCount  int             `json:"usageCount"` // 使用次数
}

// AgentCreateInput 创建代理输入
type AgentCreateInput struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Type        string          `json:"type" validate:"required"`
	Category    string          `json:"category" validate:"required"`
	Icon        string          `json:"icon"`
	Definition  json.RawMessage `json:"definition" validate:"required"`
	IsPublic    bool            `json:"isPublic"`
}

// AgentUpdateInput 更新代理输入
type AgentUpdateInput struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Category    string          `json:"category"`
	Icon        string          `json:"icon"`
	Definition  json.RawMessage `json:"definition"`
	IsPublic    bool            `json:"isPublic"`
}

// CreateAgent 创建新代理
func CreateAgent(db *sql.DB, userID *int64, input AgentCreateInput) (*Agent, error) {
	query := `
		INSERT INTO agents (name, description, type, category, icon, definition, is_public, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var agent Agent
	agent.Name = input.Name
	agent.Description = input.Description
	agent.Type = input.Type
	agent.Category = input.Category
	agent.Icon = input.Icon
	agent.Definition = input.Definition
	agent.IsPublic = input.IsPublic
	agent.UserID = userID
	agent.UsageCount = 0

	var definitionBytes []byte
	if len(input.Definition) > 0 {
		definitionBytes = []byte(input.Definition)
	} else {
		definitionBytes = []byte("{}")
	}

	err := db.QueryRow(
		query,
		agent.Name,
		agent.Description,
		agent.Type,
		agent.Category,
		agent.Icon,
		definitionBytes,
		agent.IsPublic,
		agent.UserID,
	).Scan(&agent.ID, &agent.CreatedAt, &agent.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("创建代理失败: %w", err)
	}

	return &agent, nil
}

// GetAgent 获取代理
func GetAgent(db *sql.DB, id int64) (*Agent, error) {
	query := `
		SELECT id, name, description, type, category, icon, definition, is_public, user_id, created_at, updated_at, 
		       (SELECT COUNT(*) FROM workflow_executions WHERE agent_id = agents.id) as usage_count
		FROM agents
		WHERE id = ?
	`

	var agent Agent
	var definitionBytes []byte
	var userID sql.NullInt64

	err := db.QueryRow(query, id).Scan(
		&agent.ID,
		&agent.Name,
		&agent.Description,
		&agent.Type,
		&agent.Category,
		&agent.Icon,
		&definitionBytes,
		&agent.IsPublic,
		&userID,
		&agent.CreatedAt,
		&agent.UpdatedAt,
		&agent.UsageCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("代理不存在")
		}
		return nil, fmt.Errorf("获取代理失败: %w", err)
	}

	if userID.Valid {
		agent.UserID = &userID.Int64
	} else {
		agent.UserID = nil
	}

	agent.Definition = json.RawMessage(definitionBytes)

	return &agent, nil
}

// ListAgents 获取代理列表
func ListAgents(db *sql.DB, userID *int64, category string, isPublic *bool, limit, offset int) ([]*Agent, error) {
	query := `
		SELECT id, name, description, type, category, icon, definition, is_public, user_id, created_at, updated_at,
		       (SELECT COUNT(*) FROM workflow_executions WHERE agent_id = agents.id) as usage_count
		FROM agents
		WHERE 1=1
	`
	args := []interface{}{}

	// 添加筛选条件
	if userID != nil {
		query += " AND (user_id = ? OR is_public = TRUE)"
		args = append(args, *userID)
	}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	if isPublic != nil {
		query += " AND is_public = ?"
		args = append(args, *isPublic)
	}

	// 添加排序和分页
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询代理列表失败: %w", err)
	}
	defer rows.Close()

	var agents []*Agent
	for rows.Next() {
		var agent Agent
		var definitionBytes []byte
		var userIDNull sql.NullInt64

		err := rows.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Description,
			&agent.Type,
			&agent.Category,
			&agent.Icon,
			&definitionBytes,
			&agent.IsPublic,
			&userIDNull,
			&agent.CreatedAt,
			&agent.UpdatedAt,
			&agent.UsageCount,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描代理数据失败: %w", err)
		}

		agent.Definition = json.RawMessage(definitionBytes)

		if userIDNull.Valid {
			agent.UserID = &userIDNull.Int64
		}

		agents = append(agents, &agent)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历代理数据失败: %w", err)
	}

	return agents, nil
}

// UpdateAgent 更新代理
func (a *Agent) Update(db *sql.DB, input AgentUpdateInput) error {
	query := `
		UPDATE agents
		SET 
			name = COALESCE(?, name),
			description = COALESCE(?, description),
			type = COALESCE(?, type),
			category = COALESCE(?, category),
			icon = COALESCE(?, icon),
			definition = COALESCE(?, definition),
			is_public = COALESCE(?, is_public),
			updated_at = NOW()
		WHERE id = ?
		RETURNING updated_at
	`

	var definitionBytes []byte
	if len(input.Definition) > 0 {
		definitionBytes = []byte(input.Definition)
	}

	err := db.QueryRow(
		query,
		nullString(input.Name),
		nullString(input.Description),
		nullString(input.Type),
		nullString(input.Category),
		nullString(input.Icon),
		definitionBytes,
		input.IsPublic,
		a.ID,
	).Scan(&a.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("代理不存在")
		}
		return fmt.Errorf("更新代理失败: %w", err)
	}

	// 更新内存中的代理对象
	if input.Name != "" {
		a.Name = input.Name
	}
	if input.Description != "" {
		a.Description = input.Description
	}
	if input.Type != "" {
		a.Type = input.Type
	}
	if input.Category != "" {
		a.Category = input.Category
	}
	if input.Icon != "" {
		a.Icon = input.Icon
	}
	if len(input.Definition) > 0 {
		a.Definition = input.Definition
	}
	a.IsPublic = input.IsPublic

	return nil
}

// DeleteAgent 删除代理
func (a *Agent) Delete(db *sql.DB) error {
	query := "DELETE FROM agents WHERE id = ?"

	result, err := db.Exec(query, a.ID)
	if err != nil {
		return fmt.Errorf("删除代理失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取受影响行数失败: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("代理不存在")
	}

	return nil
}

// GetAgentCategories 获取所有代理分类
func GetAgentCategories(db *sql.DB) ([]string, error) {
	query := "SELECT DISTINCT category FROM agents ORDER BY category"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询代理分类失败: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("扫描分类数据失败: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历分类数据失败: %w", err)
	}

	return categories, nil
}

// 辅助函数：将空字符串转为SQL NULL
func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
