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
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户身份")
			return
		}

		// 解析查询参数
		query := r.URL.Query()

		// 工作流ID筛选（可选）
		var workflowID *int64
		if wfID := query.Get("workflow_id"); wfID != "" {
			id, err := strconv.ParseInt(wfID, 10, 64)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "无效的工作流ID")
				return
			}
			workflowID = &id
		}

		// 状态筛选（可选）
		var status *models.ExecutionStatus
		if statusStr := query.Get("status"); statusStr != "" {
			s := models.ExecutionStatus(statusStr)
			status = &s
		}

		// 分页参数
		limit := 20
		offset := 0

		if limitStr := query.Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
				if limit > 100 {
					limit = 100 // 限制最大查询数量
				}
			}
		}

		if offsetStr := query.Get("offset"); offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}

		// 查询执行历史
		executions, err := models.ListExecutions(db.DB, workflowID, userID, status, limit, offset)
		if err != nil {
			logger.Errorw("获取执行历史失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取执行历史失败")
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
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户身份")
			return
		}

		// 获取路径参数
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的执行记录ID")
			return
		}

		// 获取执行记录
		execution, err := models.GetExecution(db.DB, id)
		if err != nil {
			logger.Errorw("获取执行记录失败", "id", id, "error", err)
			utils.RespondWithError(w, http.StatusNotFound, "执行记录不存在")
			return
		}

		// 验证用户权限
		if execution.UserID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "无权访问此执行记录")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, execution)
	}
}

// CancelExecution 取消执行中的工作流
func CancelExecution(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户身份")
			return
		}

		// 获取路径参数
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的执行记录ID")
			return
		}

		// 取消执行
		err = models.CancelExecution(db.DB, id, userID)
		if err != nil {
			logger.Errorw("取消执行失败", "id", id, "error", err)

			// 根据错误类型返回不同的状态码
			switch err.Error() {
			case "执行记录不存在":
				utils.RespondWithError(w, http.StatusNotFound, err.Error())
			case "无权访问此执行记录":
				utils.RespondWithError(w, http.StatusForbidden, err.Error())
			case "只能取消处于等待或运行状态的执行":
				utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			default:
				utils.RespondWithError(w, http.StatusInternalServerError, "取消执行失败")
			}
			return
		}

		// 返回成功消息
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "执行已成功取消",
		})
	}
}
