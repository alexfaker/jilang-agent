package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexfaker/jilang-agent/config"
	"github.com/alexfaker/jilang-agent/models"
	"github.com/alexfaker/jilang-agent/pkg/database"
	"github.com/alexfaker/jilang-agent/utils"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var validate = validator.New()

// Login 处理用户登录
func Login(db *database.DB, logger *zap.SugaredLogger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.UserLoginInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的请求数据")
			return
		}

		if err := validate.Struct(input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "验证失败: "+err.Error())
			return
		}

		// 查找用户
		user, err := models.GetUserByUsername(db.DB, input.Username)
		if err != nil {
			logger.Errorw("查找用户失败", "username", input.Username, "error", err)
			utils.RespondWithError(w, http.StatusUnauthorized, "用户名或密码不正确")
			return
		}

		// 验证密码
		if !user.CheckPassword(input.Password) {
			logger.Warnw("密码验证失败", "username", input.Username)
			utils.RespondWithError(w, http.StatusUnauthorized, "用户名或密码不正确")
			return
		}

		// 生成令牌
		token, err := utils.GenerateJWT(user.ID, user.Username, user.Role, cfg.Auth.JWTSecret, cfg.Auth.TokenExpiration)
		if err != nil {
			logger.Errorw("生成JWT失败", "userID", user.ID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "登录时发生错误")
			return
		}

		// 更新最后登录时间
		if err := user.UpdateLastLogin(db.DB); err != nil {
			logger.Warnw("更新最后登录时间失败", "userID", user.ID, "error", err)
			// 不中断处理，继续返回令牌
		}

		// 返回令牌
		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":          user.ID,
				"username":    user.Username,
				"email":       user.Email,
				"fullName":    user.FullName,
				"avatar":      user.Avatar,
				"role":        user.Role,
				"createdAt":   user.CreatedAt,
				"lastLoginAt": user.LastLoginAt,
			},
		})
	}
}

// Register 处理用户注册
func Register(db *database.DB, logger *zap.SugaredLogger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.UserRegisterInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "无效的请求数据")
			return
		}

		if err := validate.Struct(input); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "验证失败: "+err.Error())
			return
		}

		// 创建用户（包含了检查用户名和邮箱是否已存在的逻辑）
		user, err := models.CreateUser(db.DB, input)
		if err != nil {
			logger.Errorw("创建用户失败", "username", input.Username, "error", err)

			// 检查错误类型并返回适当的状态码
			if err.Error() == "用户名已存在" || err.Error() == "邮箱已存在" {
				utils.RespondWithError(w, http.StatusConflict, err.Error())
			} else {
				utils.RespondWithError(w, http.StatusInternalServerError, "注册时发生错误")
			}
			return
		}

		// 生成令牌
		token, err := utils.GenerateJWT(user.ID, user.Username, user.Role, cfg.Auth.JWTSecret, cfg.Auth.TokenExpiration)
		if err != nil {
			logger.Errorw("生成JWT失败", "userID", user.ID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "注册时发生错误")
			return
		}

		// 更新最后登录时间
		if err := user.UpdateLastLogin(db.DB); err != nil {
			logger.Warnw("更新最后登录时间失败", "userID", user.ID, "error", err)
			// 不中断处理，继续返回令牌
		}

		// 返回用户信息和令牌
		utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":          user.ID,
				"username":    user.Username,
				"email":       user.Email,
				"fullName":    user.FullName,
				"avatar":      user.Avatar,
				"role":        user.Role,
				"createdAt":   user.CreatedAt,
				"lastLoginAt": user.LastLoginAt,
			},
		})
	}
}

// RefreshToken 刷新JWT令牌
func RefreshToken(db *database.DB, logger *zap.SugaredLogger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求上下文中获取用户信息（在认证中间件中设置）
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的令牌")
			return
		}

		username, ok := r.Context().Value("username").(string)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的令牌")
			return
		}

		role, ok := r.Context().Value("role").(string)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的令牌")
			return
		}

		// 检查用户是否存在
		user, err := models.GetUserByID(db.DB, userID)
		if err != nil {
			logger.Errorw("刷新令牌时获取用户失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户")
			return
		}

		// 生成新令牌
		token, err := utils.GenerateJWT(userID, username, role, cfg.Auth.JWTSecret, cfg.Auth.TokenExpiration)
		if err != nil {
			logger.Errorw("刷新JWT失败", "userID", userID, "error", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "刷新令牌时发生错误")
			return
		}

		// 返回新令牌和用户信息
		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":          user.ID,
				"username":    user.Username,
				"email":       user.Email,
				"fullName":    user.FullName,
				"avatar":      user.Avatar,
				"role":        user.Role,
				"createdAt":   user.CreatedAt,
				"lastLoginAt": user.LastLoginAt,
			},
		})
	}
}

// Logout 用户登出
func Logout(db *database.DB, logger *zap.SugaredLogger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置清除cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   cfg.Server.Secure,
			SameSite: http.SameSiteStrictMode,
		})

		// 返回成功信息
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "已成功登出",
		})
	}
}
