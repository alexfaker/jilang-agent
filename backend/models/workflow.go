package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	AgentID      *int64          `json:"agentId"`
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

// UpdateExecutionStatus 更新执行状态
func UpdateExecutionStatus(db *sql.DB, id int64, status ExecutionStatus, logs string, errorMessage string, outputData json.RawMessage) error {
	var query string
	var args []interface{}

	if status == ExecutionStatusSuccess || status == ExecutionStatusFailed || status == ExecutionStatusCancelled {
		// 如果是完成状态，则设置完成时间和持续时间
		now := time.Now()

		query = `
			UPDATE workflow_executions
			SET 
				status = ?,
				completed_at = ?,
				duration = TIMESTAMPDIFF(SECOND, started_at, ?),
				logs = ?,
				error_message = ?,
				output_data = ?
			WHERE id = ?
		`
		args = []interface{}{
			string(status),
			now,
			now,
			logs,
			nullString(errorMessage),
			nullBytes([]byte(outputData)),
			id,
		}
	} else {
		// 如果是中间状态，则只更新状态和日志
		query = `
			UPDATE workflow_executions
			SET 
				status = ?,
				logs = ?
			WHERE id = ?
		`
		args = []interface{}{
			string(status),
			logs,
			id,
		}
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("更新执行状态失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取受影响行数失败: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("执行记录不存在")
	}

	return nil
}

// 辅助函数：将空字符串转为SQL NULL
func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// 辅助函数：将空字节切片转为SQL NULL
func nullBytes(b []byte) interface{} {
	if len(b) == 0 {
		return nil
	}
	return b
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
