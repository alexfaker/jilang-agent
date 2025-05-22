package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// WorkflowStatus 工作流状态
type WorkflowStatus string

const (
	WorkflowStatusDraft    WorkflowStatus = "draft"
	WorkflowStatusActive   WorkflowStatus = "active"
	WorkflowStatusInactive WorkflowStatus = "inactive"
	WorkflowStatusArchived WorkflowStatus = "archived"
)

// Workflow 工作流模型
type Workflow struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	UserID      int64           `json:"userId"`
	Status      WorkflowStatus  `json:"status"`
	Definition  json.RawMessage `json:"definition"` // JSON格式的工作流定义
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	LastRunAt   *time.Time      `json:"lastRunAt"`
	RunCount    int             `json:"runCount"`
}

// ExecutionStatus 执行状态
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusSuccess   ExecutionStatus = "success"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
)

// WorkflowExecution 工作流执行记录
type WorkflowExecution struct {
	ID           int64           `json:"id"`
	WorkflowID   int64           `json:"workflowId"`
	UserID       int64           `json:"userId"`
	Status       ExecutionStatus `json:"status"`
	StartedAt    time.Time       `json:"startedAt"`
	CompletedAt  *time.Time      `json:"completedAt"`
	Duration     int             `json:"duration"` // 执行时长（秒）
	Logs         string          `json:"logs"`
	ErrorMessage string          `json:"errorMessage"`
	InputData    json.RawMessage `json:"inputData"`  // 输入数据
	OutputData   json.RawMessage `json:"outputData"` // 输出数据
}

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
}

// WorkflowCreateInput 创建工作流输入
type WorkflowCreateInput struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition" validate:"required"`
	Status      WorkflowStatus  `json:"status" validate:"oneof=draft active"`
}

// WorkflowUpdateInput 更新工作流输入
type WorkflowUpdateInput struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition"`
	Status      WorkflowStatus  `json:"status" validate:"omitempty,oneof=draft active inactive archived"`
}

// ExecutionCreateInput 创建执行记录输入
type ExecutionCreateInput struct {
	WorkflowID int64           `json:"workflowId" validate:"required"`
	InputData  json.RawMessage `json:"inputData"`
}

// CreateWorkflow 创建新工作流
func CreateWorkflow(db *sql.DB, userID int64, input WorkflowCreateInput) (*Workflow, error) {
	// TODO: 实现工作流创建逻辑
	return nil, nil
}

// GetWorkflow 获取工作流
func GetWorkflow(db *sql.DB, id int64) (*Workflow, error) {
	// TODO: 实现获取工作流逻辑
	return nil, nil
}

// ListWorkflows 获取用户的工作流列表
func ListWorkflows(db *sql.DB, userID int64, status *WorkflowStatus, limit, offset int) ([]*Workflow, error) {
	// TODO: 实现工作流列表查询
	return nil, nil
}

// UpdateWorkflow 更新工作流
func (w *Workflow) Update(db *sql.DB, input WorkflowUpdateInput) error {
	// TODO: 实现工作流更新逻辑
	return nil
}

// DeleteWorkflow 删除工作流
func (w *Workflow) Delete(db *sql.DB) error {
	// TODO: 实现工作流删除逻辑
	return nil
}

// ExecuteWorkflow 执行工作流
func (w *Workflow) Execute(db *sql.DB, userID int64, input json.RawMessage) (*WorkflowExecution, error) {
	// TODO: 实现工作流执行逻辑
	return nil, nil
}

// ListExecutions 获取工作流执行历史
func ListExecutions(db *sql.DB, workflowID *int64, userID int64, status *ExecutionStatus, limit, offset int) ([]*WorkflowExecution, error) {
	// TODO: 实现执行历史查询
	return nil, nil
}

// GetExecution 获取单个执行记录
func GetExecution(db *sql.DB, id int64) (*WorkflowExecution, error) {
	// TODO: 实现获取执行记录逻辑
	return nil, nil
}
