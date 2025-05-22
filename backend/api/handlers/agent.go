package handlers

import (
	"net/http"
	"strconv"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/alexfaker/jilang-agent/utils"
	"go.uber.org/zap"
)

// ListAgents 获取可用的代理列表
func ListAgents(w http.ResponseWriter, r *http.Request) {
	// 获取分页参数
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // 默认值
	offset := 0 // 默认值

	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 20
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}
	}

	// 从数据库中获取代理列表
	agents, total, err := models.ListAgents(database.DB, limit, offset)
	if err != nil {
		zap.L().Error("获取代理列表失败", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取代理列表失败")
		return
	}

	// 返回代理列表
	utils.RespondWithPagination(w, http.StatusOK, total, limit, offset, agents)
}

// GetAgent 根据ID获取代理详情
func GetAgent(w http.ResponseWriter, r *http.Request) {
	// 获取代理ID参数
	agentIDStr := r.URL.Query().Get("id")
	if agentIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "缺少代理ID")
		return
	}

	agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "无效的代理ID")
		return
	}

	// 从数据库中获取代理详情
	agent, err := models.GetAgentByID(database.DB, agentID)
	if err != nil {
		zap.L().Error("获取代理详情失败", zap.Error(err), zap.Int64("agentID", agentID))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取代理详情失败")
		return
	}

	// 如果找不到代理
	if agent == nil {
		utils.RespondWithError(w, http.StatusNotFound, "代理不存在")
		return
	}

	// 返回代理详情
	utils.RespondWithJSON(w, http.StatusOK, agent)
}
