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

// ListAgents 获取代理列表
func ListAgents(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户身份")
			return
		}

		// 解析查询参数
		query := r.URL.Query()

		// 类别筛选（可选）
		category := query.Get("category")

		// 是否公开筛选（可选）
		var isPublic *bool
		if publicStr := query.Get("is_public"); publicStr != "" {
			if publicStr == "true" {
				trueValue := true
				isPublic = &trueValue
			} else if publicStr == "false" {
				falseValue := false
				isPublic = &falseValue
			}
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

		// 查询代理列表
		agents, err := models.ListAgents(db.DB, &userID, category, isPublic, limit, offset)
		if err != nil {
			logger.Errorw("获取代理列表失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取代理列表失败")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, agents)
	}
}

// GetAgent 获取单个代理
func GetAgent(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
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
			utils.RespondWithError(w, http.StatusBadRequest, "无效的代理ID")
			return
		}

		// 获取代理
		agent, err := models.GetAgent(db.DB, id)
		if err != nil {
			logger.Errorw("获取代理失败", "id", id, "error", err)
			utils.RespondWithError(w, http.StatusNotFound, "代理不存在")
			return
		}

		// 验证用户权限（只能查看公开的代理或自己的代理）
		if !agent.IsPublic && (agent.UserID == nil || *agent.UserID != userID) {
			utils.RespondWithError(w, http.StatusForbidden, "无权访问此代理")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, agent)
	}
}

// GetAgentCategories 获取代理分类列表
func GetAgentCategories(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取所有代理分类
		categories, err := models.GetAgentCategories(db.DB)
		if err != nil {
			logger.Errorw("获取代理分类失败", "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "获取代理分类失败")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, categories)
	}
}

// CreateAgent 创建代理
func CreateAgent(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户ID
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户身份")
			return
		}

		// 解析请求体
		var input models.AgentCreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的请求数据")
			return
		}

		// 验证输入
		if err := validate.Struct(input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "验证失败: "+err.Error())
			return
		}

		// 创建代理
		agent, err := models.CreateAgent(db.DB, &userID, input)
		if err != nil {
			logger.Errorw("创建代理失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "创建代理失败")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusCreated, agent)
	}
}

// UpdateAgent 更新代理
func UpdateAgent(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
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
			utils.RespondWithError(w, http.StatusBadRequest, "无效的代理ID")
			return
		}

		// 获取代理
		agent, err := models.GetAgent(db.DB, id)
		if err != nil {
			logger.Errorw("获取代理失败", "id", id, "error", err)
			utils.RespondWithError(w, http.StatusNotFound, "代理不存在")
			return
		}

		// 验证用户权限（只能更新自己的代理）
		if agent.UserID == nil || *agent.UserID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "无权更新此代理")
			return
		}

		// 解析请求体
		var input models.AgentUpdateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的请求数据")
			return
		}

		// 更新代理
		if err := agent.Update(db.DB, input); err != nil {
			logger.Errorw("更新代理失败", "id", id, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "更新代理失败")
			return
		}

		// 返回结果
		utils.RespondWithJSON(w, http.StatusOK, agent)
	}
}

// DeleteAgent 删除代理
func DeleteAgent(db *database.DB, logger *zap.SugaredLogger) http.HandlerFunc {
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
			utils.RespondWithError(w, http.StatusBadRequest, "无效的代理ID")
			return
		}

		// 获取代理
		agent, err := models.GetAgent(db.DB, id)
		if err != nil {
			logger.Errorw("获取代理失败", "id", id, "error", err)
			utils.RespondWithError(w, http.StatusNotFound, "代理不存在")
			return
		}

		// 验证用户权限（只能删除自己的代理）
		if agent.UserID == nil || *agent.UserID != userID {
			utils.RespondWithError(w, http.StatusForbidden, "无权删除此代理")
			return
		}

		// 删除代理
		if err := agent.Delete(db.DB); err != nil {
			logger.Errorw("删除代理失败", "id", id, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "删除代理失败")
			return
		}

		// 返回成功消息
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "代理已成功删除",
		})
	}
}
