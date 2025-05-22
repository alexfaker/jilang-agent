package handlers

import (
	"net/http"
	"strconv"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/alexfaker/jilang-agent/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// ListExecutions 获取执行历史列表
func ListExecutions(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "未认证")
			return
		}

		// 获取查询参数
		query := r.URL.Query()

		// 工作流ID筛选
		var workflowIDFilter *int64
		if workflowIDStr := query.Get("workflow_id"); workflowIDStr != "" {
			if wID, err := strconv.ParseInt(workflowIDStr, 10, 64); err == nil {
				workflowIDFilter = &wID
			}
		}

		// 状态筛选
		var statusFilter *models.ExecutionStatus
		if status := query.Get("status"); status != "" {
			s := models.ExecutionStatus(status)
			statusFilter = &s
		}

		// 分页参数
		limit := 20
		offset := 0
		if limitStr := query.Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}
		if offsetStr := query.Get("offset"); offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}

		// 查询执行历史列表
		executions, err := models.ListExecutions(db.DB, workflowIDFilter, userID, statusFilter, limit, offset)
		if err != nil {
			logger.Errorw("获取执行历史列表失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取执行历史列表时发生错误")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, executions)
	}
}

// GetExecution 获取单个执行记录
func GetExecution(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "未认证")
			return
		}

		// 获取执行ID
		executionID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的执行ID")
			return
		}

		// 查询执行记录
		execution, err := models.GetExecution(db.DB, executionID)
		if err != nil {
			logger.Errorw("获取执行记录失败", "executionID", executionID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取执行记录时发生错误")
			return
		}

		if execution == nil {
			utils.RespondWithError(w, http.StatusNotFound, "执行记录不存在")
			return
		}

		// 权限检查
		if execution.UserID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "无权访问此执行记录")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, execution)
	}
}
