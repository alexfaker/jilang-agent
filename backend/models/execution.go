package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
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

	// 关联信息（非数据库字段，用于API返回）
	WorkflowName string `json:"workflowName,omitempty"`
	AgentName    string `json:"agentName,omitempty"`
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
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Running   int `json:"running"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
	Canceled  int `json:"canceled"`
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
func CreateExecution(db *sql.DB, userID int64, input ExecutionCreateInput) (*WorkflowExecution, error) {
	// 检查工作流是否存在并且属于当前用户
	var workflowUserID int64
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

// GetExecution 获取单个执行记录
func GetExecution(db *sql.DB, id int64) (*WorkflowExecution, error) {
	query := `
		SELECT e.id, e.workflow_id, e.user_id, e.agent_id, e.status, e.started_at, e.completed_at, 
		       e.duration, e.logs, e.error_message, e.input_data, e.output_data,
		       w.name as workflow_name, a.name as agent_name
		FROM workflow_executions e
		JOIN workflows w ON e.workflow_id = w.id
		LEFT JOIN agents a ON e.agent_id = a.id
		WHERE e.id = ?
	`

	var execution WorkflowExecution
	var inputDataBytes, outputDataBytes []byte
	var completedAt sql.NullTime
	var agentID sql.NullInt64
	var logs, errorMessage sql.NullString
	var agentName sql.NullString

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
		&execution.WorkflowName,
		&agentName,
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
	if agentName.Valid {
		execution.AgentName = agentName.String
	}

	execution.InputData = json.RawMessage(inputDataBytes)
	execution.OutputData = json.RawMessage(outputDataBytes)

	return &execution, nil
}

// ListExecutions 获取工作流执行历史
func ListExecutions(db *sql.DB, workflowID *int64, userID int64, status *ExecutionStatus, limit, offset int) ([]*WorkflowExecution, error) {
	query := `
		SELECT e.id, e.workflow_id, e.user_id, e.agent_id, e.status, e.started_at, e.completed_at, 
		       e.duration, e.logs, e.error_message, e.input_data, e.output_data,
		       w.name as workflow_name, a.name as agent_name
		FROM workflow_executions e
		JOIN workflows w ON e.workflow_id = w.id
		LEFT JOIN agents a ON e.agent_id = a.id
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
		var agentName sql.NullString

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
			&execution.WorkflowName,
			&agentName,
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
		if agentName.Valid {
			execution.AgentName = agentName.String
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
				logs = CONCAT(COALESCE(logs, ''), ?),
				error_message = ?,
				output_data = ?
			WHERE id = ?
		`
		args = []interface{}{
			string(status),
			now,
			now,
			"\n" + logs, // 附加日志而不是替换
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
				logs = CONCAT(COALESCE(logs, ''), ?)
			WHERE id = ?
		`
		args = []interface{}{
			string(status),
			"\n" + logs, // 附加日志而不是替换
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

// CancelExecution 取消执行
func CancelExecution(db *sql.DB, id int64, userID int64) error {
	// 验证执行记录所属
	var execUserID int64
	var status string
	err := db.QueryRow(`
		SELECT e.user_id, e.status 
		FROM workflow_executions e 
		WHERE e.id = ?
	`, id).Scan(&execUserID, &status)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("执行记录不存在")
		}
		return fmt.Errorf("查询执行记录失败: %w", err)
	}

	if execUserID != userID {
		return errors.New("无权访问此执行记录")
	}

	// 只能取消处于等待或运行状态的执行
	if status != string(ExecutionStatusPending) && status != string(ExecutionStatusRunning) {
		return errors.New("只能取消处于等待或运行状态的执行")
	}

	// 更新状态为已取消
	return UpdateExecutionStatus(db, id, ExecutionStatusCancelled, "用户取消执行", "用户手动取消", nil)
}

// GetExecutionStats 获取执行统计数据
func GetExecutionStats(db *sql.DB, userID int64, workflowID *int64, days int) (*ExecutionStatsData, error) {
	stats := &ExecutionStatsData{
		StatusCounts: make(map[string]int),
		DateCounts:   make(map[string]int),
	}

	// 设置默认查询时间范围
	if days <= 0 {
		days = 30 // 默认查询最近30天
	}

	// 基础查询条件
	baseQuery := `
		FROM workflow_executions e
		JOIN workflows w ON e.workflow_id = w.id
		WHERE w.user_id = ? AND e.started_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
	`
	args := []interface{}{userID, days}

	// 添加工作流ID筛选（如果提供）
	if workflowID != nil {
		baseQuery += " AND e.workflow_id = ?"
		args = append(args, *workflowID)
	}

	// 查询总执行次数
	err := db.QueryRow("SELECT COUNT(*) "+baseQuery, args...).Scan(&stats.TotalCount)
	if err != nil {
		return nil, fmt.Errorf("查询总执行次数失败: %w", err)
	}

	// 查询成功执行次数
	successArgs := append([]interface{}{string(ExecutionStatusSuccess)}, args...)
	err = db.QueryRow("SELECT COUNT(*) "+baseQuery+" AND e.status = ?", successArgs...).Scan(&stats.SuccessCount)
	if err != nil {
		return nil, fmt.Errorf("查询成功执行次数失败: %w", err)
	}

	// 查询失败执行次数
	failedArgs := append([]interface{}{string(ExecutionStatusFailed)}, args...)
	err = db.QueryRow("SELECT COUNT(*) "+baseQuery+" AND e.status = ?", failedArgs...).Scan(&stats.FailedCount)
	if err != nil {
		return nil, fmt.Errorf("查询失败执行次数失败: %w", err)
	}

	// 计算成功率
	if stats.TotalCount > 0 {
		stats.SuccessRate = float64(stats.SuccessCount) / float64(stats.TotalCount) * 100
	}

	// 查询平均执行时间（只考虑已完成的执行）
	avgDurationQuery := `
		SELECT COALESCE(AVG(duration), 0)
		` + baseQuery + `
		AND e.completed_at IS NOT NULL
	`
	err = db.QueryRow(avgDurationQuery, args...).Scan(&stats.AverageDuration)
	if err != nil {
		return nil, fmt.Errorf("查询平均执行时间失败: %w", err)
	}

	// 查询各状态的执行次数
	statusQuery := `
		SELECT e.status, COUNT(*)
		` + baseQuery + `
		GROUP BY e.status
	`
	statusRows, err := db.Query(statusQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("查询状态统计失败: %w", err)
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var status string
		var count int
		if err := statusRows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("扫描状态统计数据失败: %w", err)
		}
		stats.StatusCounts[status] = count
	}

	if err := statusRows.Err(); err != nil {
		return nil, fmt.Errorf("遍历状态统计数据失败: %w", err)
	}

	// 查询按日期的执行次数
	dateQuery := `
		SELECT DATE(e.started_at) as date, COUNT(*)
		` + baseQuery + `
		GROUP BY DATE(e.started_at)
		ORDER BY date
	`
	dateRows, err := db.Query(dateQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("查询日期统计失败: %w", err)
	}
	defer dateRows.Close()

	for dateRows.Next() {
		var date string
		var count int
		if err := dateRows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("扫描日期统计数据失败: %w", err)
		}
		stats.DateCounts[date] = count
	}

	if err := dateRows.Err(); err != nil {
		return nil, fmt.Errorf("遍历日期统计数据失败: %w", err)
	}

	return stats, nil
}

// GetExecutionStats 获取执行记录统计数据
func GetExecutionStats(db *sql.DB, userID int64) (*ExecutionStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN status = 'running' THEN 1 ELSE 0 END) as running,
			SUM(CASE WHEN status = 'succeeded' THEN 1 ELSE 0 END) as succeeded,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed,
			SUM(CASE WHEN status = 'canceled' THEN 1 ELSE 0 END) as canceled
		FROM workflow_executions
		WHERE user_id = ?
	`

	var stats ExecutionStats
	err := db.QueryRow(query, userID).Scan(
		&stats.Total,
		&stats.Pending,
		&stats.Running,
		&stats.Succeeded,
		&stats.Failed,
		&stats.Canceled,
	)
	if err != nil {
		return nil, fmt.Errorf("获取执行统计失败: %w", err)
	}

	return &stats, nil
}

// GetExecutionStatsByWorkflow 获取按工作流分组的执行统计数据
func GetExecutionStatsByWorkflow(db *sql.DB, userID int64) ([]*WorkflowExecutionStats, error) {
	query := `
		SELECT 
			e.workflow_id,
			w.name as workflow_name,
			COUNT(*) as total,
			SUM(CASE WHEN e.status = 'succeeded' THEN 1 ELSE 0 END) as succeeded,
			SUM(CASE WHEN e.status = 'failed' THEN 1 ELSE 0 END) as failed,
			MAX(e.created_at) as last_executed
		FROM workflow_executions e
		JOIN workflows w ON e.workflow_id = w.id
		WHERE e.user_id = ?
		GROUP BY e.workflow_id, w.name
		ORDER BY total DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询工作流执行统计失败: %w", err)
	}
	defer rows.Close()

	var stats []*WorkflowExecutionStats
	for rows.Next() {
		var stat WorkflowExecutionStats
		var lastExecuted sql.NullTime

		err := rows.Scan(
			&stat.WorkflowID,
			&stat.WorkflowName,
			&stat.Total,
			&stat.Succeeded,
			&stat.Failed,
			&lastExecuted,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描工作流执行统计数据失败: %w", err)
		}

		if lastExecuted.Valid {
			stat.LastExecuted = &lastExecuted.Time
		}

		// 计算成功率
		if stat.Total > 0 {
			stat.SuccessRate = float64(stat.Succeeded) / float64(stat.Total) * 100
		}

		stats = append(stats, &stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历工作流执行统计数据失败: %w", err)
	}

	return stats, nil
}

// GetExecutionStatsByDateRange 获取指定日期范围内的执行统计数据
func GetExecutionStatsByDateRange(db *sql.DB, userID int64, startDate, endDate time.Time) (*ExecutionStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN status = 'running' THEN 1 ELSE 0 END) as running,
			SUM(CASE WHEN status = 'succeeded' THEN 1 ELSE 0 END) as succeeded,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed,
			SUM(CASE WHEN status = 'canceled' THEN 1 ELSE 0 END) as canceled
		FROM workflow_executions
		WHERE user_id = ? AND created_at >= ? AND created_at < ?
	`

	var stats ExecutionStats
	err := db.QueryRow(query, userID, startDate, endDate).Scan(
		&stats.Total,
		&stats.Pending,
		&stats.Running,
		&stats.Succeeded,
		&stats.Failed,
		&stats.Canceled,
	)
	if err != nil {
		return nil, fmt.Errorf("获取日期范围内执行统计失败: %w", err)
	}

	return &stats, nil
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
