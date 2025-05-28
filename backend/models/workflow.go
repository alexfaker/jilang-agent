package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// WorkflowStatus 工作流状态
type WorkflowStatus string

const (
	WorkflowStatusDraft    WorkflowStatus = "draft"
	WorkflowStatusActive   WorkflowStatus = "active"
	WorkflowStatusInactive WorkflowStatus = "inactive"
	WorkflowStatusArchived WorkflowStatus = "archived"
)

// Workflow 工作流模型 - 用户购买的工作流实例
type Workflow struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string          `json:"name" gorm:"type:varchar(100);not null"`
	Description string          `json:"description" gorm:"type:text"`
	UserID      int64           `json:"userId" gorm:"column:user_id;index;not null"`
	AgentID     *int64          `json:"agentId" gorm:"column:agent_id;index"` // 关联的代理ID（购买来源）
	Status      WorkflowStatus  `json:"status" gorm:"type:varchar(20);default:'draft';not null"`
	Definition  json.RawMessage `json:"definition" gorm:"type:json"`            // JSON格式的工作流定义
	PurchasedAt *time.Time      `json:"purchasedAt" gorm:"column:purchased_at"` // 购买时间
	CreatedAt   time.Time       `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	LastRunAt   *time.Time      `json:"lastRunAt" gorm:"column:last_run_at"`
	RunCount    int             `json:"runCount" gorm:"column:run_count;default:0"`
}

// TableName 指定表名
func (Workflow) TableName() string {
	return "workflows"
}

// WorkflowAgent 工作流中使用的代理模型
type WorkflowAgent struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string          `json:"name" gorm:"type:varchar(100);not null"`
	Description string          `json:"description" gorm:"type:text"`
	Type        string          `json:"type" gorm:"type:varchar(50);not null"`
	Category    string          `json:"category" gorm:"type:varchar(50);index"`
	Icon        string          `json:"icon" gorm:"type:varchar(255)"`
	Definition  json.RawMessage `json:"definition" gorm:"type:json"` // 代理定义JSON
	IsPublic    bool            `json:"isPublic" gorm:"column:is_public;default:false"`
	UserID      *int64          `json:"userId" gorm:"column:user_id;index"` // 为空表示系统预置
	CreatedAt   time.Time       `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (WorkflowAgent) TableName() string {
	return "agents"
}

// WorkflowCreateInput 创建工作流输入
type WorkflowCreateInput struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition" validate:"required"`
	Status      WorkflowStatus  `json:"status" validate:"oneof=draft active"`
	AgentID     *int64          `json:"agentId"` // 关联的代理ID（购买来源）
}

// WorkflowUpdateInput 更新工作流输入
type WorkflowUpdateInput struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition"`
	Status      WorkflowStatus  `json:"status" validate:"omitempty,oneof=draft active inactive archived"`
}

// CreateWorkflow 使用GORM创建新工作流
func CreateWorkflow(db *gorm.DB, userID int64, input WorkflowCreateInput) (*Workflow, error) {
	// 验证工作流状态是否有效
	if input.Status == "" {
		input.Status = WorkflowStatusDraft
	}

	// 创建工作流对象
	workflow := &Workflow{
		Name:        input.Name,
		Description: input.Description,
		UserID:      userID,
		AgentID:     input.AgentID,
		Status:      input.Status,
		Definition:  input.Definition,
		RunCount:    0,
	}

	// 如果是从代理购买的，设置购买时间
	if input.AgentID != nil {
		now := time.Now()
		workflow.PurchasedAt = &now
	}

	// 保存到数据库
	if err := db.Create(workflow).Error; err != nil {
		return nil, err
	}

	return workflow, nil
}

// GetWorkflow 使用GORM获取工作流
func GetWorkflow(db *gorm.DB, id int64) (*Workflow, error) {
	var workflow Workflow
	if err := db.First(&workflow, id).Error; err != nil {
		return nil, err
	}
	return &workflow, nil
}

// ListWorkflows 使用GORM获取用户的工作流列表
func ListWorkflows(db *gorm.DB, userID int64, status *WorkflowStatus, limit, offset int) ([]*Workflow, error) {
	var workflows []*Workflow
	query := db.Where("user_id = ?", userID)

	// 添加状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 执行查询
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&workflows).Error
	if err != nil {
		return nil, err
	}

	return workflows, nil
}

// Update 使用GORM更新工作流
func (w *Workflow) Update(db *gorm.DB, input WorkflowUpdateInput) error {
	// 准备更新数据
	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if len(input.Definition) > 0 {
		updates["definition"] = input.Definition
	}

	// 更新工作流
	return db.Model(w).Updates(updates).Error
}

// Delete 使用GORM删除工作流
func (w *Workflow) Delete(db *gorm.DB) error {
	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的执行记录
		if err := tx.Where("workflow_id = ?", w.ID).Delete(&WorkflowExecution{}).Error; err != nil {
			return err
		}

		// 删除工作流
		if err := tx.Delete(w).Error; err != nil {
			return err
		}

		return nil
	})
}

// Execute 使用GORM执行工作流
func (w *Workflow) Execute(db *gorm.DB, userID int64, inputData json.RawMessage) (*WorkflowExecution, error) {
	// 验证工作流状态
	if w.Status != WorkflowStatusActive {
		return nil, gorm.ErrInvalidTransaction
	}

	// 创建执行记录
	now := time.Now()
	execution := WorkflowExecution{
		WorkflowID: w.ID,
		UserID:     userID,
		Status:     ExecutionStatusPending,
		StartedAt:  now,
		InputData:  inputData,
	}

	// 保存到数据库
	if err := db.Create(&execution).Error; err != nil {
		return nil, err
	}

	// 更新工作流的最后运行时间和运行次数
	updates := map[string]interface{}{
		"last_run_at": now,
		"run_count":   gorm.Expr("run_count + ?", 1),
	}

	if err := db.Model(w).Updates(updates).Error; err != nil {
		// 这里不返回错误，因为执行记录已经创建成功
		// 但应该记录日志
	} else {
		// 更新内存中的工作流对象
		w.LastRunAt = &now
		w.RunCount++
	}

	return &execution, nil
}

// CountWorkflows 使用GORM统计用户的工作流数量
func CountWorkflows(db *gorm.DB, userID int64) (int64, error) {
	var count int64
	err := db.Model(&Workflow{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CreateWorkflowGorm 使用GORM创建新工作流
func CreateWorkflowGorm(db *gorm.DB, userID int64, input WorkflowCreateInput) (*Workflow, error) {
	// 验证工作流状态是否有效
	if input.Status == "" {
		input.Status = WorkflowStatusDraft
	} else if input.Status != WorkflowStatusDraft && input.Status != WorkflowStatusActive {
		return nil, errors.New("无效的工作流状态")
	}

	// 验证工作流定义是否是有效的JSON
	if len(input.Definition) == 0 {
		return nil, errors.New("工作流定义不能为空")
	}

	var js json.RawMessage
	if err := json.Unmarshal(input.Definition, &js); err != nil {
		return nil, fmt.Errorf("工作流定义不是有效的JSON: %w", err)
	}

	// 创建工作流对象
	workflow := &Workflow{
		Name:        input.Name,
		Description: input.Description,
		UserID:      userID,
		Status:      input.Status,
		Definition:  input.Definition,
		RunCount:    0,
	}

	// 保存到数据库
	if err := db.Create(workflow).Error; err != nil {
		return nil, fmt.Errorf("创建工作流失败: %w", err)
	}

	return workflow, nil
}

// GetWorkflowGorm 使用GORM获取工作流
func GetWorkflowGorm(db *gorm.DB, id int64) (*Workflow, error) {
	var workflow Workflow
	if err := db.First(&workflow, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工作流不存在")
		}
		return nil, fmt.Errorf("获取工作流失败: %w", err)
	}
	return &workflow, nil
}

// ListWorkflowsGorm 使用GORM获取用户的工作流列表
func ListWorkflowsGorm(db *gorm.DB, userID int64, status *WorkflowStatus, limit, offset int) ([]*Workflow, error) {
	var workflows []*Workflow
	query := db.Where("user_id = ?", userID)

	// 添加状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 执行查询
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&workflows).Error
	if err != nil {
		return nil, fmt.Errorf("查询工作流列表失败: %w", err)
	}

	return workflows, nil
}

// UpdateGorm 使用GORM更新工作流
func (w *Workflow) UpdateGorm(db *gorm.DB, input WorkflowUpdateInput) error {
	// 验证工作流状态是否有效
	if input.Status != "" &&
		input.Status != WorkflowStatusDraft &&
		input.Status != WorkflowStatusActive &&
		input.Status != WorkflowStatusInactive &&
		input.Status != WorkflowStatusArchived {
		return errors.New("无效的工作流状态")
	}

	// 验证工作流定义是否是有效的JSON
	if len(input.Definition) > 0 {
		var js json.RawMessage
		if err := json.Unmarshal(input.Definition, &js); err != nil {
			return fmt.Errorf("工作流定义不是有效的JSON: %w", err)
		}
	}

	// 准备更新数据
	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if len(input.Definition) > 0 {
		updates["definition"] = input.Definition
	}

	// 更新工作流
	if err := db.Model(w).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新工作流失败: %w", err)
	}

	return nil
}

// DeleteGorm 使用GORM删除工作流
func (w *Workflow) DeleteGorm(db *gorm.DB) error {
	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的执行记录
		if err := tx.Where("workflow_id = ?", w.ID).Delete(&WorkflowExecution{}).Error; err != nil {
			return fmt.Errorf("删除执行记录失败: %w", err)
		}

		// 删除工作流
		if err := tx.Delete(w).Error; err != nil {
			return fmt.Errorf("删除工作流失败: %w", err)
		}

		return nil
	})
}

// ExecuteGorm 使用GORM执行工作流
func (w *Workflow) ExecuteGorm(db *gorm.DB, userID int64, inputData json.RawMessage) (*WorkflowExecution, error) {
	// 验证工作流状态
	if w.Status != WorkflowStatusActive {
		return nil, errors.New("只能执行处于活动状态的工作流")
	}

	// 创建执行记录
	now := time.Now()
	execution := WorkflowExecution{
		WorkflowID: w.ID,
		UserID:     userID,
		Status:     ExecutionStatusPending,
		StartedAt:  now,
		InputData:  inputData,
	}

	// 保存到数据库
	if err := db.Create(&execution).Error; err != nil {
		return nil, fmt.Errorf("创建执行记录失败: %w", err)
	}

	// 更新工作流的最后运行时间和运行次数
	updates := map[string]interface{}{
		"last_run_at": now,
		"run_count":   gorm.Expr("run_count + ?", 1),
	}

	if err := db.Model(w).Updates(updates).Error; err != nil {
		// 这里不返回错误，因为执行记录已经创建成功
		fmt.Printf("警告: 更新工作流最后运行时间和运行次数失败: %v\n", err)
	} else {
		// 更新内存中的工作流对象
		w.LastRunAt = &now
		w.RunCount++
	}

	return &execution, nil
}

// CountWorkflowsGorm 使用GORM统计用户的工作流数量
func CountWorkflowsGorm(db *gorm.DB, userID int64) (int64, error) {
	var count int64
	if err := db.Model(&Workflow{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计工作流数量失败: %w", err)
	}
	return count, nil
}
