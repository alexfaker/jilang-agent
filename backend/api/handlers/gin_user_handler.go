package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GinUserHandler 处理用户相关的请求
type GinUserHandler struct {
	DB        *gorm.DB
	Logger    *zap.Logger
	Validator *validator.Validate
}

// NewGinUserHandler 创建一个新的GinUserHandler实例
func NewGinUserHandler(db *gorm.DB, logger *zap.Logger) *GinUserHandler {
	return &GinUserHandler{
		DB:        db,
		Logger:    logger,
		Validator: validator.New(),
	}
}

// GinUserProfileResponse 用户资料响应结构
type GinUserProfileResponse struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

// GetUserProfile 获取用户资料
func (h *GinUserHandler) GetUserProfile(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}
	uid := userID.(int64)

	// 获取用户信息
	var user models.User
	if err := h.DB.First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户失败", zap.Error(err), zap.Int64("user_id", uid))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户失败: " + err.Error(),
			})
		}
		return
	}

	// 移除敏感信息
	profile := GinUserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   profile,
	})
}

// GetUserByID 根据ID获取用户
func (h *GinUserHandler) GetUserByID(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的用户ID",
		})
		return
	}

	// 获取用户信息
	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户失败", zap.Error(err), zap.Int64("user_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户失败: " + err.Error(),
			})
		}
		return
	}

	// 移除敏感信息
	profile := GinUserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   profile,
	})
}

// GinUpdateProfileRequest 更新用户资料请求结构
type GinUpdateProfileRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserProfile 更新用户资料
func (h *GinUserHandler) UpdateUserProfile(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}
	uid := userID.(int64)

	// 解析请求体
	var req GinUpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 验证输入
	if err := h.Validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "验证失败: " + err.Error(),
		})
		return
	}

	// 获取当前用户信息
	var user models.User
	if err := h.DB.First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户失败", zap.Error(err), zap.Int64("user_id", uid))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户失败: " + err.Error(),
			})
		}
		return
	}

	// 更新用户信息
	updates := map[string]interface{}{
		"email": req.Email,
	}

	if err := h.DB.Model(&user).Updates(updates).Error; err != nil {
		h.Logger.Error("更新用户失败", zap.Error(err), zap.Int64("user_id", uid))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新用户失败: " + err.Error(),
		})
		return
	}

	// 重新获取更新后的用户信息
	if err := h.DB.First(&user, uid).Error; err != nil {
		h.Logger.Error("获取更新后的用户失败", zap.Error(err), zap.Int64("user_id", uid))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取更新后的用户失败: " + err.Error(),
		})
		return
	}

	// 移除敏感信息
	profile := GinUserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	// 返回更新后的用户信息
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "用户资料更新成功",
		"data":    profile,
	})
}

// GetCurrentUser 获取当前用户信息
func (h *GinUserHandler) GetCurrentUser(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}
	uid := userID.(int64)

	// 获取用户信息
	var user models.User
	if err := h.DB.First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户失败", zap.Error(err), zap.Int64("user_id", uid))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户失败: " + err.Error(),
			})
		}
		return
	}

	// 移除敏感信息
	profile := GinUserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   profile,
	})
}

// UpdateUser 更新用户信息
func (h *GinUserHandler) UpdateUser(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}
	uid := userID.(int64)

	// 解析请求体
	var req GinUpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 验证输入
	if err := h.Validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "验证失败: " + err.Error(),
		})
		return
	}

	// 更新用户信息
	updates := map[string]interface{}{
		"email": req.Email,
	}

	if err := h.DB.Model(&models.User{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
		h.Logger.Error("更新用户失败", zap.Error(err), zap.Int64("user_id", uid))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新用户失败: " + err.Error(),
		})
		return
	}

	// 获取更新后的用户信息
	var user models.User
	if err := h.DB.First(&user, uid).Error; err != nil {
		h.Logger.Error("获取更新后的用户失败", zap.Error(err), zap.Int64("user_id", uid))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取更新后的用户失败: " + err.Error(),
		})
		return
	}

	// 移除敏感信息
	profile := GinUserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	// 返回更新后的用户信息
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "用户资料更新成功",
		"data":    profile,
	})
}

// GinChangePasswordRequest 修改密码请求结构
type GinChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword 修改密码
func (h *GinUserHandler) ChangePassword(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}
	uid := userID.(int64)

	// 解析请求体
	var req GinChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 获取用户信息
	var user models.User
	if err := h.DB.First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户失败", zap.Error(err), zap.Int64("user_id", uid))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户失败: " + err.Error(),
			})
		}
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "旧密码不正确",
		})
		return
	}

	// 验证新密码格式
	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "新密码长度不能少于8个字符",
		})
		return
	}

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("生成密码哈希失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "修改密码失败: " + err.Error(),
		})
		return
	}

	// 更新密码
	if err := h.DB.Model(&user).Update("password_hash", string(hashedPassword)).Error; err != nil {
		h.Logger.Error("更新密码失败", zap.Error(err), zap.Int64("user_id", uid))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "修改密码失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "密码修改成功",
	})
}
