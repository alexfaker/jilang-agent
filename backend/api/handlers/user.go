package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/alexfaker/jilang-agent/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// GetUserProfile 获取用户个人资料
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(int64)

	// 获取用户信息
	user, err := models.GetUserByID(database.DB, userID)
	if err != nil {
		zap.L().Error("获取用户信息失败", zap.Error(err), zap.Int64("userID", userID))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	// 移除敏感信息
	user.PasswordHash = ""

	// 返回用户信息
	utils.RespondWithJSON(w, http.StatusOK, user)
}

// GetUserByID 根据ID获取用户
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取用户ID
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 获取用户信息
	user, err := models.GetUserByID(database.DB, userID)
	if err != nil {
		zap.L().Error("获取用户信息失败", zap.Error(err), zap.Int64("userID", userID))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	// 移除敏感信息
	user.PasswordHash = ""

	// 返回用户信息
	utils.RespondWithJSON(w, http.StatusOK, user)
}

// UpdateUserProfile 更新用户个人资料
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(int64)

	// 解析请求体
	var input models.UpdateUserInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}
	defer r.Body.Close()

	// 验证输入
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "数据验证失败: "+err.Error())
		return
	}

	// 获取当前用户信息
	user, err := models.GetUserByID(database.DB, userID)
	if err != nil {
		zap.L().Error("获取用户信息失败", zap.Error(err), zap.Int64("userID", userID))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	// 更新用户信息
	updatedUser, err := models.UpdateUser(database.DB, userID, input)
	if err != nil {
		zap.L().Error("更新用户信息失败", zap.Error(err), zap.Int64("userID", userID))
		utils.RespondWithError(w, http.StatusInternalServerError, "更新用户信息失败")
		return
	}

	// 移除敏感信息
	updatedUser.PasswordHash = ""

	// 返回更新后的用户信息
	utils.RespondWithJSON(w, http.StatusOK, updatedUser)
}

// ChangePassword 修改用户密码
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	// 从上下文获取用户ID
	userID := r.Context().Value("userID").(int64)

	// 解析请求体
	var input models.ChangePasswordInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}
	defer r.Body.Close()

	// 验证输入
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "数据验证失败: "+err.Error())
		return
	}

	// 获取用户信息
	user, err := models.GetUserByID(database.DB, userID)
	if err != nil {
		zap.L().Error("获取用户信息失败", zap.Error(err), zap.Int64("userID", userID))
		utils.RespondWithError(w, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	// 验证当前密码
	if !models.CheckPasswordHash(input.CurrentPassword, user.PasswordHash) {
		utils.RespondWithError(w, http.StatusBadRequest, "当前密码不正确")
		return
	}

	// 修改密码
	err = models.ChangePassword(database.DB, userID, input.NewPassword)
	if err != nil {
		zap.L().Error("修改密码失败", zap.Error(err), zap.Int64("userID", userID))
		utils.RespondWithError(w, http.StatusInternalServerError, "修改密码失败")
		return
	}

	// 返回成功消息
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "密码修改成功"})
}
