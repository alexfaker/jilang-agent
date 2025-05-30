package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

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
	ID           int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	WorkflowID   int64           `json:"workflowId" gorm:"column:workflow_id;index;not null"`
	UserID       string          `json:"userID" gorm:"column:user_id;index;not null"`
	AgentID      *int64          `json:"agentId" gorm:"column:agent_id;index"`
	Status       ExecutionStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	StartedAt    time.Time       `json:"startedAt" gorm:"column:started_at;not null"`
	CompletedAt  *time.Time      `json:"completedAt" gorm:"column:completed_at"`
	Duration     int             `json:"duration" gorm:"default:0"` // 执行时长（秒）
	Logs         string          `json:"logs" gorm:"type:text"`
	ErrorMessage string          `json:"errorMessage" gorm:"column:error_message;type:text"`
	InputData    json.RawMessage `json:"inputData" gorm:"column:input_data;type:json"`   // 输入数据
	OutputData   json.RawMessage `json:"outputData" gorm:"column:output_data;type:json"` // 输出数据
	CreatedAt    time.Time       `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time       `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`

	// GORM 关联关系
	Workflow *Workflow `json:"workflow,omitempty" gorm:"foreignKey:WorkflowID;references:ID"`

	// 关联信息（非数据库字段，用于API返回）
	WorkflowName string `json:"workflowName,omitempty" gorm:"-"`
	AgentName    string `json:"agentName,omitempty" gorm:"-"`
}

// TableName 指定表名
func (WorkflowExecution) TableName() string {
	return "workflow_executions"
}

// ExecutionCreateInput 创建执行记录输入
type ExecutionCreateInput struct {
	WorkflowID int64           `json:"workflowId" validate:"required"`
	AgentID    *int64          `json:"agentId"`
	InputData  json.RawMessage `json:"inputData"`
}

// ExecutionStatsData 执行统计数据
type ExecutionStatsData struct {
	TotalCount      int            `json:"totalCount"`
	SuccessCount    int            `json:"successCount"`
	FailedCount     int            `json:"failedCount"`
	SuccessRate     float64        `json:"successRate"`     // 百分比
	AverageDuration float64        `json:"averageDuration"` // 秒
	StatusCounts    map[string]int `json:"statusCounts"`
	DateCounts      map[string]int `json:"dateCounts"` // 按日期统计
}

// ExecutionStats 执行记录统计数据
type ExecutionStats struct {
	Total     int64 `json:"total"`
	Pending   int   `json:"pending"`
	Running   int   `json:"running"`
	Succeeded int   `json:"succeeded"`
	Failed    int   `json:"failed"`
	Canceled  int   `json:"canceled"`
}

// WorkflowExecutionStats 工作流执行统计数据
type WorkflowExecutionStats struct {
	WorkflowID   int64      `json:"workflowId"`
	WorkflowName string     `json:"workflowName"`
	Total        int        `json:"total"`
	Succeeded    int        `json:"succeeded"`
	Failed       int        `json:"failed"`
	SuccessRate  float64    `json:"successRate"`
	LastExecuted *time.Time `json:"lastExecuted"`
}

// CreateExecution 创建新执行记录
func CreateExecution(db *sql.DB, userID string, input ExecutionCreateInput) (*WorkflowExecution, error) {
	// 检查工作流是否存在并且属于当前用户
	var workflowUserID string
	var workflowStatus string
	err := db.QueryRow("SELECT user_id, status FROM workflows WHERE id = ?", input.WorkflowID).Scan(&workflowUserID, &workflowStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("工作流不存在")
		}
		return nil, fmt.Errorf("获取工作流信息失败: %w", err)
	}

	if workflowUserID != userID {
		return nil, errors.New("无权访问此工作流")
	}

	if workflowStatus != string(WorkflowStatusActive) {
		return nil, errors.New("只能执行处于活动状态的工作流")
	}

	// 如果指定了代理，检查代理是否存在
	if input.AgentID != nil {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM agents WHERE id = ?", *input.AgentID).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("检查代理失败: %w", err)
		}
		if count == 0 {
			return nil, errors.New("代理不存在")
		}
	}

	// 创建执行记录
	now := time.Now()
	execution := WorkflowExecution{
		WorkflowID: input.WorkflowID,
		UserID:     userID,
		AgentID:    input.AgentID,
		Status:     ExecutionStatusPending,
		StartedAt:  now,
		InputData:  input.InputData,
	}

	// 插入执行记录
	query := `
		INSERT INTO workflow_executions (workflow_id, user_id, agent_id, status, started_at, input_data)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`

	err = db.QueryRow(
		query,
		execution.WorkflowID,
		execution.UserID,
		execution.AgentID,
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

	_, err = db.Exec(updateQuery, now, input.WorkflowID)
	if err != nil {
		// 记录警告，但不返回错误
		fmt.Printf("警告: 更新工作流最后运行时间和运行次数失败: %v\n", err)
	}

	return &execution, nil
}

// CreateExecutionGorm 使用GORM创建新执行记录
func CreateExecutionGorm(db *gorm.DB, userID string, input ExecutionCreateInput) (*WorkflowExecution, error) {
	// 检查工作流是否存在并且属于当前用户
	var workflow Workflow
	if err := db.Select("user_id, status").Where("id = ?", input.WorkflowID).First(&workflow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工作流不存在")
		}
		return nil, fmt.Errorf("获取工作流信息失败: %w", err)
	}

	// 将userID转换为字符串进行比较
	userIDStr := fmt.Sprintf("%d", userID)
	if workflow.UserID != userIDStr {
		return nil, errors.New("无权访问此工作流")
	}

	if workflow.Status != WorkflowStatusActive {
		return nil, errors.New("只能执行处于活动状态的工作流")
	}

	// 如果指定了代理，检查代理是否存在
	if input.AgentID != nil {
		var count int64
		if err := db.Model(&Agent{}).Where("id = ?", *input.AgentID).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("检查代理失败: %w", err)
		}
		if count == 0 {
			return nil, errors.New("代理不存在")
		}
	}

	// 创建执行记录
	now := time.Now()
	execution := WorkflowExecution{
		WorkflowID: input.WorkflowID,
		UserID:     userID,
		AgentID:    input.AgentID,
		Status:     ExecutionStatusPending,
		StartedAt:  now,
		InputData:  input.InputData,
	}

	// 保存到数据库
	if err := db.Create(&execution).Error; err != nil {
		return nil, fmt.Errorf("创建执行记录失败: %w", err)
	}

	// 更新工作流的最后运行时间和运行次数
	if err := db.Model(&Workflow{}).Where("id = ?", input.WorkflowID).Updates(map[string]interface{}{
		"last_run_at": now,
		"run_count":   gorm.Expr("run_count + ?", 1),
	}).Error; err != nil {
		// 记录警告，但不返回错误
		fmt.Printf("警告: 更新工作流最后运行时间和运行次数失败: %v\n", err)
	}

	return &execution, nil
}

// GetExecutionGorm 使用GORM获取单个执行记录
func GetExecutionGorm(db *gorm.DB, id int64) (*WorkflowExecution, error) {
	var execution WorkflowExecution

	// 使用关联查询获取执行记录及相关信息
	err := db.Table("workflow_executions as e").
		Select("e.*, w.name as workflow_name, a.name as agent_name").
		Joins("JOIN workflows w ON e.workflow_id = w.id").
		Joins("LEFT JOIN agents a ON e.agent_id = a.id").
		Where("e.id = ?", id).
		Scan(&execution).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("执行记录不存在")
		}
		return nil, fmt.Errorf("获取执行记录失败: %w", err)
	}

	return &execution, nil
}

// ListExecutionsGorm 使用GORM获取执行记录列表
func ListExecutionsGorm(db *gorm.DB, workflowID *int64, userID int64, status *ExecutionStatus, limit, offset int) ([]*WorkflowExecution, error) {
	var executions []*WorkflowExecution

	// 构建查询
	query := db.Table("workflow_executions as e").
		Select("e.*, w.name as workflow_name, a.name as agent_name").
		Joins("JOIN workflows w ON e.workflow_id = w.id").
		Joins("LEFT JOIN agents a ON e.agent_id = a.id").
		Where("w.user_id = ?", userID)

	// 添加工作流ID筛选
	if workflowID != nil {
		query = query.Where("e.workflow_id = ?", *workflowID)
	}

	// 添加状态筛选
	if status != nil {
		query = query.Where("e.status = ?", *status)
	}

	// 执行查询
	err := query.Order("e.started_at DESC").Limit(limit).Offset(offset).Scan(&executions).Error
	if err != nil {
		return nil, fmt.Errorf("查询执行记录列表失败: %w", err)
	}

	return executions, nil
}

// UpdateExecutionStatusGorm 使用GORM更新执行记录状态
func UpdateExecutionStatusGorm(db *gorm.DB, id int64, status ExecutionStatus, logs string, errorMessage string, outputData json.RawMessage) error {
	// 准备更新数据
	updates := map[string]interface{}{
		"status": status,
	}

	// 如果状态为成功或失败，设置完成时间和持续时间
	if status == ExecutionStatusSuccess || status == ExecutionStatusFailed || status == ExecutionStatusCancelled {
		now := time.Now()
		updates["completed_at"] = now

		// 获取开始时间以计算持续时间
		var startedAt time.Time
		if err := db.Model(&WorkflowExecution{}).Where("id = ?", id).Select("started_at").Scan(&startedAt).Error; err != nil {
			return fmt.Errorf("获取执行记录开始时间失败: %w", err)
		}

		// 计算持续时间（秒）
		duration := int(now.Sub(startedAt).Seconds())
		updates["duration"] = duration
	}

	// 设置日志和错误信息（如果有）
	if logs != "" {
		updates["logs"] = logs
	}
	if errorMessage != "" {
		updates["error_message"] = errorMessage
	}
	if len(outputData) > 0 {
		updates["output_data"] = outputData
	}

	// 执行更新
	if err := db.Model(&WorkflowExecution{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新执行记录状态失败: %w", err)
	}

	return nil
}

// CancelExecutionGorm 使用GORM取消执行
func CancelExecutionGorm(db *gorm.DB, id int64, userID int64) error {
	// 获取执行记录
	var execution WorkflowExecution
	if err := db.Where("id = ?", id).First(&execution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("执行记录不存在")
		}
		return fmt.Errorf("获取执行记录失败: %w", err)
	}

	// 验证用户权限
	var workflow Workflow
	if err := db.Where("id = ?", execution.WorkflowID).Select("user_id").First(&workflow).Error; err != nil {
		return fmt.Errorf("获取工作流信息失败: %w", err)
	}

	// 将userID转换为字符串进行比较
	userIDStr := fmt.Sprintf("%d", userID)
	if workflow.UserID != userIDStr {
		return errors.New("无权访问此执行记录")
	}

	// 验证执行状态
	if execution.Status != ExecutionStatusPending && execution.Status != ExecutionStatusRunning {
		return errors.New("只能取消处于等待或运行状态的执行")
	}

	// 更新状态为已取消
	return UpdateExecutionStatusGorm(db, id, ExecutionStatusCancelled, "执行已被用户取消", "", nil)
}

// GetExecutionStatsGorm 使用GORM获取执行统计数据
func GetExecutionStatsGorm(db *gorm.DB, userID int64) (*ExecutionStats, error) {
	var stats ExecutionStats

	// 获取总数
	if err := db.Model(&WorkflowExecution{}).
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ?", userID).
		Count(&stats.Total).Error; err != nil {
		return nil, fmt.Errorf("统计执行记录总数失败: %w", err)
	}

	// 获取各状态数量
	statusCounts := []struct {
		Status string
		Count  int
	}{}

	if err := db.Model(&WorkflowExecution{}).
		Select("status, count(*) as count").
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ?", userID).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("统计执行记录状态分布失败: %w", err)
	}

	// 填充各状态数量
	for _, sc := range statusCounts {
		switch ExecutionStatus(sc.Status) {
		case ExecutionStatusPending:
			stats.Pending = sc.Count
		case ExecutionStatusRunning:
			stats.Running = sc.Count
		case ExecutionStatusSuccess:
			stats.Succeeded = sc.Count
		case ExecutionStatusFailed:
			stats.Failed = sc.Count
		case ExecutionStatusCancelled:
			stats.Canceled = sc.Count
		}
	}

	return &stats, nil
}

// GetExecutionStatsByWorkflowGorm 使用GORM获取按工作流分组的执行统计数据
func GetExecutionStatsByWorkflowGorm(db *gorm.DB, userID int64) ([]*WorkflowExecutionStats, error) {
	var stats []*WorkflowExecutionStats

	// 查询每个工作流的执行统计
	rows, err := db.Raw(`
		SELECT 
			w.id as workflow_id, 
			w.name as workflow_name,
			COUNT(*) as total,
			SUM(CASE WHEN e.status = 'success' THEN 1 ELSE 0 END) as succeeded,
			SUM(CASE WHEN e.status = 'failed' THEN 1 ELSE 0 END) as failed,
			MAX(e.started_at) as last_executed
		FROM workflow_executions e
		JOIN workflows w ON e.workflow_id = w.id
		WHERE w.user_id = ?
		GROUP BY w.id, w.name
		ORDER BY total DESC
	`, userID).Rows()
	if err != nil {
		return nil, fmt.Errorf("查询工作流执行统计失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stat WorkflowExecutionStats
		var lastExecuted sql.NullTime
		var total, succeeded, failed int

		if err := rows.Scan(&stat.WorkflowID, &stat.WorkflowName, &total, &succeeded, &failed, &lastExecuted); err != nil {
			return nil, fmt.Errorf("扫描工作流执行统计数据失败: %w", err)
		}

		stat.Total = total
		stat.Succeeded = succeeded
		stat.Failed = failed

		// 计算成功率
		if total > 0 {
			stat.SuccessRate = float64(succeeded) / float64(total) * 100
		}

		if lastExecuted.Valid {
			stat.LastExecuted = &lastExecuted.Time
		}

		stats = append(stats, &stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历工作流执行统计数据失败: %w", err)
	}

	return stats, nil
}

// GetExecutionStatsByDateRangeGorm 使用GORM获取指定日期范围内的执行统计数据
func GetExecutionStatsByDateRangeGorm(db *gorm.DB, userID int64, startDate, endDate time.Time) (*ExecutionStats, error) {
	var stats ExecutionStats

	// 获取总数
	if err := db.Model(&WorkflowExecution{}).
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ? AND workflow_executions.started_at BETWEEN ? AND ?", userID, startDate, endDate).
		Count(&stats.Total).Error; err != nil {
		return nil, fmt.Errorf("统计日期范围内执行记录总数失败: %w", err)
	}

	// 获取各状态数量
	statusCounts := []struct {
		Status string
		Count  int
	}{}

	if err := db.Model(&WorkflowExecution{}).
		Select("status, count(*) as count").
		Joins("JOIN workflows ON workflow_executions.workflow_id = workflows.id").
		Where("workflows.user_id = ? AND workflow_executions.started_at BETWEEN ? AND ?", userID, startDate, endDate).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("统计日期范围内执行记录状态分布失败: %w", err)
	}

	// 填充各状态数量
	for _, sc := range statusCounts {
		switch ExecutionStatus(sc.Status) {
		case ExecutionStatusPending:
			stats.Pending = sc.Count
		case ExecutionStatusRunning:
			stats.Running = sc.Count
		case ExecutionStatusSuccess:
			stats.Succeeded = sc.Count
		case ExecutionStatusFailed:
			stats.Failed = sc.Count
		case ExecutionStatusCancelled:
			stats.Canceled = sc.Count
		}
	}

	return &stats, nil
}

// 辅助函数：将空字节数组转为SQL NULL
func nullBytes(b []byte) interface{} {
	if len(b) == 0 {
		return nil
	}
	return b
}
