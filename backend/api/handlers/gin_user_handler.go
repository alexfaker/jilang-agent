package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
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
	UserID      string     `json:"userID"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FullName    string     `json:"fullName"`
	Avatar      string     `json:"avatar"`
	Bio         string     `json:"bio"`
	Timezone    string     `json:"timezone"`
	Language    string     `json:"language"`
	Theme       string     `json:"theme"`
	Role        string     `json:"role"`
	Points      int        `json:"points"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
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
	uid := userID.(string)

	// 根据UserID获取用户信息
	var user models.User
	if err := h.DB.Where("user_id = ?", uid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户失败", zap.Error(err), zap.String("user_id", uid))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户失败: " + err.Error(),
			})
		}
		return
	}

	// 构建完整的用户资料响应
	profile := GinUserProfileResponse{
		ID:          user.ID,
		UserID:      user.UserID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		Timezone:    user.Timezone,
		Language:    user.Language,
		Theme:       user.Theme,
		Role:        user.Role,
		Points:      user.Points,
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
	Email    string `json:"email" validate:"omitempty,email"`
	FullName string `json:"fullName"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Timezone string `json:"timezone"`
	Language string `json:"language"`
	Theme    string `json:"theme"`
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
	uid := userID.(string)

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

	// 根据UserID获取用户信息
	var user models.User
	if err := h.DB.Where("user_id = ?", uid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户失败", zap.Error(err), zap.String("user_id", uid))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户失败: " + err.Error(),
			})
		}
		return
	}

	// 使用User模型的Update方法
	input := models.UserUpdateInput{
		Email:    req.Email,
		FullName: req.FullName,
		Avatar:   req.Avatar,
		Bio:      req.Bio,
		Timezone: req.Timezone,
		Language: req.Language,
		Theme:    req.Theme,
	}

	if err := user.Update(h.DB, input); err != nil {
		if err == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "邮箱已被使用",
			})
			return
		}
		h.Logger.Error("更新用户失败", zap.Error(err), zap.String("user_id", uid))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新用户失败: " + err.Error(),
		})
		return
	}

	// 构建完整的用户资料响应
	profile := GinUserProfileResponse{
		ID:          user.ID,
		UserID:      user.UserID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		Timezone:    user.Timezone,
		Language:    user.Language,
		Theme:       user.Theme,
		Role:        user.Role,
		Points:      user.Points,
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
	uid := userID.(string)

	// 根据UserID获取用户信息
	var user models.User
	if err := h.DB.Where("user_id = ?", uid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户失败", zap.Error(err), zap.String("user_id", uid))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户失败: " + err.Error(),
			})
		}
		return
	}

	// 构建完整的用户资料响应
	profile := GinUserProfileResponse{
		ID:          user.ID,
		UserID:      user.UserID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		Timezone:    user.Timezone,
		Language:    user.Language,
		Theme:       user.Theme,
		Role:        user.Role,
		Points:      user.Points,
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
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
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
	uid := userID.(string)

	// 解析请求体
	var req GinChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 根据UserID获取用户信息
	var user models.User
	if err := h.DB.Where("user_id = ?", uid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("获取用户失败", zap.Error(err), zap.String("user_id", uid))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "获取用户失败: " + err.Error(),
			})
		}
		return
	}

	// 使用User模型的ChangePassword方法
	input := models.PasswordChangeInput{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	if err := user.ChangePassword(h.DB, input); err != nil {
		if err == gorm.ErrInvalidValue {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "当前密码不正确",
			})
			return
		}
		h.Logger.Error("修改密码失败", zap.Error(err), zap.String("user_id", uid))
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

// UploadAvatar 上传头像
func (h *GinUserHandler) UploadAvatar(c *gin.Context) {
	// 从请求上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户身份",
		})
		return
	}
	uid := userID.(string)

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "获取上传文件失败: " + err.Error(),
		})
		return
	}

	// 检查文件大小（最大2MB）
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "文件大小不能超过2MB",
		})
		return
	}

	// 检查文件类型
	if file.Header.Get("Content-Type") != "image/jpeg" &&
		file.Header.Get("Content-Type") != "image/png" &&
		file.Header.Get("Content-Type") != "image/gif" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "只支持 JPG, PNG, GIF 格式的图片",
		})
		return
	}

	// 生成文件名（使用用户ID + 时间戳）
	filename := "avatar_" + uid + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"

	// 保存文件到public/uploads/avatars目录
	uploadPath := "public/uploads/avatars/" + filename
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		h.Logger.Error("保存头像文件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "保存头像文件失败: " + err.Error(),
		})
		return
	}

	// 更新用户头像URL
	avatarURL := "/uploads/avatars/" + filename
	var user models.User
	if err := h.DB.Where("user_id = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "用户不存在",
		})
		return
	}

	if err := h.DB.Model(&user).Update("avatar", avatarURL).Error; err != nil {
		h.Logger.Error("更新用户头像失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新用户头像失败: " + err.Error(),
		})
		return
	}

	// 返回更新后的用户资料
	profile := GinUserProfileResponse{
		ID:          user.ID,
		UserID:      user.UserID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Avatar:      avatarURL,
		Bio:         user.Bio,
		Timezone:    user.Timezone,
		Language:    user.Language,
		Theme:       user.Theme,
		Role:        user.Role,
		Points:      user.Points,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "头像上传成功",
		"data":    profile,
	})
}
