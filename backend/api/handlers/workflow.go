package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/alexfaker/jilang-agent/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// ListWorkflows 获取工作流列表
func ListWorkflows(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "未认证")
			return
		}

		// 获取查询参数
		query := r.URL.Query()

		// 状态筛选
		var statusFilter *models.WorkflowStatus
		if status := query.Get("status"); status != "" {
			s := models.WorkflowStatus(status)
			statusFilter = &s
		}

		// 分页参数
		limit := 10
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

		// 查询工作流列表
		workflows, err := models.ListWorkflows(db.DB, userID, statusFilter, limit, offset)
		if err != nil {
			logger.Errorw("获取工作流列表失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取工作流列表时发生错误")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, workflows)
	}
}

// GetWorkflow 获取单个工作流
func GetWorkflow(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "未认证")
			return
		}

		// 获取工作流ID
		workflowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的工作流ID")
			return
		}

		// 查询工作流
		workflow, err := models.GetWorkflow(db.DB, workflowID)
		if err != nil {
			logger.Errorw("获取工作流失败", "workflowID", workflowID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取工作流时发生错误")
			return
		}

		if workflow == nil {
			utils.RespondWithError(w, http.StatusNotFound, "工作流不存在")
			return
		}

		// 权限检查
		if workflow.UserID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "无权访问此工作流")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, workflow)
	}
}

// CreateWorkflow 创建工作流
func CreateWorkflow(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "未认证")
			return
		}

		// 解析请求数据
		var input models.WorkflowCreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的请求数据")
			return
		}

		// 验证请求数据
		if err := validate.Struct(input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "验证失败: "+err.Error())
			return
		}

		// 创建工作流
		workflow, err := models.CreateWorkflow(db.DB, userID, input)
		if err != nil {
			logger.Errorw("创建工作流失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "创建工作流时发生错误")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusCreated, workflow)
	}
}

// UpdateWorkflow 更新工作流
func UpdateWorkflow(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "未认证")
			return
		}

		// 获取工作流ID
		workflowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的工作流ID")
			return
		}

		// 解析请求数据
		var input models.WorkflowUpdateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的请求数据")
			return
		}

		// 验证请求数据
		if err := validate.Struct(input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "验证失败: "+err.Error())
			return
		}

		// 查询工作流
		workflow, err := models.GetWorkflow(db.DB, workflowID)
		if err != nil {
			logger.Errorw("获取工作流失败", "workflowID", workflowID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "更新工作流时发生错误")
			return
		}

		if workflow == nil {
			utils.RespondWithError(w, http.StatusNotFound, "工作流不存在")
			return
		}

		// 权限检查
		if workflow.UserID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "无权更新此工作流")
			return
		}

		// 更新工作流
		if err := workflow.Update(db.DB, input); err != nil {
			logger.Errorw("更新工作流失败", "workflowID", workflowID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "更新工作流时发生错误")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, workflow)
	}
}

// DeleteWorkflow 删除工作流
func DeleteWorkflow(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "未认证")
			return
		}

		// 获取工作流ID
		workflowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的工作流ID")
			return
		}

		// 查询工作流
		workflow, err := models.GetWorkflow(db.DB, workflowID)
		if err != nil {
			logger.Errorw("获取工作流失败", "workflowID", workflowID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "删除工作流时发生错误")
			return
		}

		if workflow == nil {
			utils.RespondWithError(w, http.StatusNotFound, "工作流不存在")
			return
		}

		// 权限检查
		if workflow.UserID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "无权删除此工作流")
			return
		}

		// 删除工作流
		if err := workflow.Delete(db.DB); err != nil {
			logger.Errorw("删除工作流失败", "workflowID", workflowID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "删除工作流时发生错误")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "工作流已成功删除",
		})
	}
}

// ExecuteWorkflow 执行工作流
func ExecuteWorkflow(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "未认证")
			return
		}

		// 获取工作流ID
		workflowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的工作流ID")
			return
		}

		// 解析请求数据
		var inputData json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
			if err.Error() != "EOF" { // 允许空输入
				utils.RespondWithError(w, http.StatusBadRequest, "无效的请求数据")
				return
			}
			inputData = json.RawMessage("{}")
		}

		// 查询工作流
		workflow, err := models.GetWorkflow(db.DB, workflowID)
		if err != nil {
			logger.Errorw("获取工作流失败", "workflowID", workflowID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "执行工作流时发生错误")
			return
		}

		if workflow == nil {
			utils.RespondWithError(w, http.StatusNotFound, "工作流不存在")
			return
		}

		// 检查工作流状态
		if workflow.Status != models.WorkflowStatusActive {
			utils.RespondWithError(w, http.StatusBadRequest, "只能执行处于活跃状态的工作流")
			return
		}

		// 执行工作流
		execution, err := workflow.Execute(db.DB, userID, inputData)
		if err != nil {
			logger.Errorw("执行工作流失败", "workflowID", workflowID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "执行工作流时发生错误")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusAccepted, execution)
	}
}
