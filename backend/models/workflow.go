package models

import (
	"database/sql"
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

// Workflow 工作流模型
type Workflow struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string          `json:"name" gorm:"type:varchar(100);not null"`
	Description string          `json:"description" gorm:"type:text"`
	UserID      int64           `json:"userId" gorm:"column:user_id;index;not null"`
	Status      WorkflowStatus  `json:"status" gorm:"type:varchar(20);default:'draft';not null"`
	Definition  json.RawMessage `json:"definition" gorm:"type:json"` // JSON格式的工作流定义
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
}

// WorkflowUpdateInput 更新工作流输入
type WorkflowUpdateInput struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition"`
	Status      WorkflowStatus  `json:"status" validate:"omitempty,oneof=draft active inactive archived"`
}

// CreateWorkflow 创建新工作流
func CreateWorkflow(db *sql.DB, userID int64, input WorkflowCreateInput) (*Workflow, error) {
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

	// 插入新工作流
	query := `
		INSERT INTO workflows (name, description, user_id, status, definition, created_at, updated_at, run_count)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW(), 0)
		RETURNING id, created_at, updated_at
	`

	var workflow Workflow
	workflow.Name = input.Name
	workflow.Description = input.Description
	workflow.UserID = userID
	workflow.Status = input.Status
	workflow.Definition = input.Definition
	workflow.RunCount = 0

	err := db.QueryRow(
		query,
		workflow.Name,
		workflow.Description,
		workflow.UserID,
		workflow.Status,
		[]byte(workflow.Definition),
	).Scan(&workflow.ID, &workflow.CreatedAt, &workflow.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("创建工作流失败: %w", err)
	}

	return &workflow, nil
}

// GetWorkflow 获取工作流
func GetWorkflow(db *sql.DB, id int64) (*Workflow, error) {
	query := `
		SELECT id, name, description, user_id, status, definition, created_at, updated_at, last_run_at, run_count
		FROM workflows
		WHERE id = ?
	`

	var workflow Workflow
	var definitionBytes []byte
	var lastRunAt sql.NullTime

	err := db.QueryRow(query, id).Scan(
		&workflow.ID,
		&workflow.Name,
		&workflow.Description,
		&workflow.UserID,
		&workflow.Status,
		&definitionBytes,
		&workflow.CreatedAt,
		&workflow.UpdatedAt,
		&lastRunAt,
		&workflow.RunCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("工作流不存在")
		}
		return nil, fmt.Errorf("获取工作流失败: %w", err)
	}

	if lastRunAt.Valid {
		workflow.LastRunAt = &lastRunAt.Time
	}

	workflow.Definition = json.RawMessage(definitionBytes)

	return &workflow, nil
}

// ListWorkflows 获取用户的工作流列表
func ListWorkflows(db *sql.DB, userID int64, status *WorkflowStatus, limit, offset int) ([]*Workflow, error) {
	query := `
		SELECT id, name, description, user_id, status, definition, created_at, updated_at, last_run_at, run_count
		FROM workflows
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
		return nil, fmt.Errorf("查询工作流列表失败: %w", err)
	}
	defer rows.Close()

	var workflows []*Workflow
	for rows.Next() {
		var workflow Workflow
		var definitionBytes []byte
		var lastRunAt sql.NullTime

		err := rows.Scan(
			&workflow.ID,
			&workflow.Name,
			&workflow.Description,
			&workflow.UserID,
			&workflow.Status,
			&definitionBytes,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
			&lastRunAt,
			&workflow.RunCount,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描工作流数据失败: %w", err)
		}

		workflow.Definition = json.RawMessage(definitionBytes)

		if lastRunAt.Valid {
			workflow.LastRunAt = &lastRunAt.Time
		}

		workflows = append(workflows, &workflow)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历工作流数据失败: %w", err)
	}

	return workflows, nil
}

// UpdateWorkflow 更新工作流
func (w *Workflow) Update(db *sql.DB, input WorkflowUpdateInput) error {
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

	// 更新工作流
	query := `
		UPDATE workflows
		SET 
			name = COALESCE(?, name),
			description = COALESCE(?, description),
			status = COALESCE(?, status),
			definition = COALESCE(?, definition),
			updated_at = NOW()
		WHERE id = ?
		RETURNING updated_at
	`

	var definitionBytes []byte
	if len(input.Definition) > 0 {
		definitionBytes = []byte(input.Definition)
	}

	var statusStr *string
	if input.Status != "" {
		s := string(input.Status)
		statusStr = &s
	}

	err := db.QueryRow(
		query,
		nullString(input.Name),
		nullString(input.Description),
		statusStr,
		nullBytes(definitionBytes),
		w.ID,
	).Scan(&w.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("工作流不存在")
		}
		return fmt.Errorf("更新工作流失败: %w", err)
	}

	// 更新内存中的工作流对象
	if input.Name != "" {
		w.Name = input.Name
	}
	if input.Description != "" {
		w.Description = input.Description
	}
	if input.Status != "" {
		w.Status = input.Status
	}
	if len(input.Definition) > 0 {
		w.Definition = input.Definition
	}

	return nil
}

// DeleteWorkflow 删除工作流
func (w *Workflow) Delete(db *sql.DB) error {
	// 首先检查是否有与此工作流关联的执行记录
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM workflow_executions WHERE workflow_id = ?", w.ID).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查执行记录失败: %w", err)
	}

	// 如果有执行记录，则使用事务删除工作流及其执行记录
	if count > 0 {
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("开始事务失败: %w", err)
		}

		// 删除执行记录
		_, err = tx.Exec("DELETE FROM workflow_executions WHERE workflow_id = ?", w.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("删除执行记录失败: %w", err)
		}

		// 删除工作流
		_, err = tx.Exec("DELETE FROM workflows WHERE id = ?", w.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("删除工作流失败: %w", err)
		}

		// 提交事务
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("提交事务失败: %w", err)
		}
	} else {
		// 如果没有执行记录，直接删除工作流
		result, err := db.Exec("DELETE FROM workflows WHERE id = ?", w.ID)
		if err != nil {
			return fmt.Errorf("删除工作流失败: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("获取受影响行数失败: %w", err)
		}

		if rowsAffected == 0 {
			return errors.New("工作流不存在")
		}
	}

	return nil
}

// Execute 执行工作流
func (w *Workflow) Execute(db *sql.DB, userID int64, inputData json.RawMessage) (*WorkflowExecution, error) {
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

	// 插入执行记录
	query := `
		INSERT INTO workflow_executions (workflow_id, user_id, status, started_at, input_data)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		execution.WorkflowID,
		execution.UserID,
		execution.Status,
		execution.StartedAt,
		[]byte(execution.InputData),
	).Scan(&execution.ID)

	if err != nil {
		return nil, fmt.Errorf("创建执行记录失败: %w", err)
	}

	// 更新工作流的最后运行时间和运行次数
	updateQuery := `
		UPDATE workflows
		SET last_run_at = ?, run_count = run_count + 1
		WHERE id = ?
	`

	_, err = db.Exec(updateQuery, now, w.ID)
	if err != nil {
		// 这里不返回错误，因为执行记录已经创建成功，更新工作流信息失败并不影响执行过程
		// 但应该记录日志
		fmt.Printf("警告: 更新工作流最后运行时间和运行次数失败: %v\n", err)
	} else {
		// 更新内存中的工作流对象
		w.LastRunAt = &now
		w.RunCount++
	}

	// 注意：实际的工作流执行逻辑将在单独的goroutine中异步进行
	// 这里只返回初始化的执行记录
	return &execution, nil
}

// ListExecutions 获取工作流执行历史
func ListExecutions(db *sql.DB, workflowID *int64, userID int64, status *ExecutionStatus, limit, offset int) ([]*WorkflowExecution, error) {
	query := `
		SELECT e.id, e.workflow_id, e.user_id, e.agent_id, e.status, e.started_at, e.completed_at, 
		       e.duration, e.logs, e.error_message, e.input_data, e.output_data
		FROM workflow_executions e
		JOIN workflows w ON e.workflow_id = w.id
		WHERE w.user_id = ?
	`
	args := []interface{}{userID}

	// 添加工作流ID筛选
	if workflowID != nil {
		query += " AND e.workflow_id = ?"
		args = append(args, *workflowID)
	}

	// 添加状态筛选
	if status != nil {
		query += " AND e.status = ?"
		args = append(args, string(*status))
	}

	// 添加排序和分页
	query += " ORDER BY e.started_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询执行历史失败: %w", err)
	}
	defer rows.Close()

	var executions []*WorkflowExecution
	for rows.Next() {
		var execution WorkflowExecution
		var inputDataBytes, outputDataBytes []byte
		var completedAt sql.NullTime
		var agentID sql.NullInt64
		var logs, errorMessage sql.NullString

		err := rows.Scan(
			&execution.ID,
			&execution.WorkflowID,
			&execution.UserID,
			&agentID,
			&execution.Status,
			&execution.StartedAt,
			&completedAt,
			&execution.Duration,
			&logs,
			&errorMessage,
			&inputDataBytes,
			&outputDataBytes,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描执行记录失败: %w", err)
		}

		if agentID.Valid {
			execution.AgentID = &agentID.Int64
		}
		if completedAt.Valid {
			execution.CompletedAt = &completedAt.Time
		}
		if logs.Valid {
			execution.Logs = logs.String
		}
		if errorMessage.Valid {
			execution.ErrorMessage = errorMessage.String
		}

		execution.InputData = json.RawMessage(inputDataBytes)
		execution.OutputData = json.RawMessage(outputDataBytes)

		executions = append(executions, &execution)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历执行记录失败: %w", err)
	}

	return executions, nil
}

// GetExecution 获取单个执行记录
func GetExecution(db *sql.DB, id int64) (*WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, user_id, agent_id, status, started_at, completed_at, 
		       duration, logs, error_message, input_data, output_data
		FROM workflow_executions
		WHERE id = ?
	`

	var execution WorkflowExecution
	var inputDataBytes, outputDataBytes []byte
	var completedAt sql.NullTime
	var agentID sql.NullInt64
	var logs, errorMessage sql.NullString

	err := db.QueryRow(query, id).Scan(
		&execution.ID,
		&execution.WorkflowID,
		&execution.UserID,
		&agentID,
		&execution.Status,
		&execution.StartedAt,
		&completedAt,
		&execution.Duration,
		&logs,
		&errorMessage,
		&inputDataBytes,
		&outputDataBytes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("执行记录不存在")
		}
		return nil, fmt.Errorf("获取执行记录失败: %w", err)
	}

	if agentID.Valid {
		execution.AgentID = &agentID.Int64
	}
	if completedAt.Valid {
		execution.CompletedAt = &completedAt.Time
	}
	if logs.Valid {
		execution.Logs = logs.String
	}
	if errorMessage.Valid {
		execution.ErrorMessage = errorMessage.String
	}

	execution.InputData = json.RawMessage(inputDataBytes)
	execution.OutputData = json.RawMessage(outputDataBytes)

	return &execution, nil
}

// CountWorkflows 统计用户的工作流数量
func CountWorkflows(db *sql.DB, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM workflows WHERE user_id = ?`

	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("统计工作流数量失败: %w", err)
	}

	return count, nil
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
