package handlers

import (
	"net/http"
	"strconv"

	"github.com/alexfaker/jilang-agent/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GinSettingsHandler 设置处理程序
type GinSettingsHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewGinSettingsHandler 创建新的设置处理程序
func NewGinSettingsHandler(db *gorm.DB, logger *zap.Logger) *GinSettingsHandler {
	return &GinSettingsHandler{
		DB:     db,
		Logger: logger,
	}
}

// SettingsResponse 设置响应结构
type SettingsResponse struct {
	Theme struct {
		Mode  string `json:"mode"`
		Color string `json:"color"`
	} `json:"theme"`
	Language   string `json:"language"`
	DateFormat string `json:"dateFormat"`
	TimeFormat string `json:"timeFormat"`
	Timezone   string `json:"timezone"`
	Layout     struct {
		SidebarPosition string `json:"sidebarPosition"`
		ContentWidth    string `json:"contentWidth"`
		CompactMode     bool   `json:"compactMode"`
	} `json:"layout"`
	Notifications struct {
		Email struct {
			Enabled   bool   `json:"enabled"`
			Frequency string `json:"frequency"`
			Types     struct {
				WorkflowCompleted bool `json:"workflow_completed"`
				WorkflowFailed    bool `json:"workflow_failed"`
				AgentError        bool `json:"agent_error"`
				NewUpdates        bool `json:"new_updates"`
				SecurityAlerts    bool `json:"security_alerts"`
			} `json:"types"`
		} `json:"email"`
		Browser struct {
			Enabled bool   `json:"enabled"`
			Sound   string `json:"sound"`
			Types   struct {
				WorkflowCompleted bool `json:"workflow_completed"`
				WorkflowFailed    bool `json:"workflow_failed"`
				AgentError        bool `json:"agent_error"`
				SecurityAlerts    bool `json:"security_alerts"`
			} `json:"types"`
		} `json:"browser"`
		InApp struct {
			Enabled  bool   `json:"enabled"`
			Position string `json:"position"`
			Types    struct {
				WorkflowCompleted bool `json:"workflow_completed"`
				WorkflowFailed    bool `json:"workflow_failed"`
				AgentError        bool `json:"agent_error"`
				SecurityAlerts    bool `json:"security_alerts"`
				NewUpdates        bool `json:"new_updates"`
			} `json:"types"`
		} `json:"inApp"`
	} `json:"notifications"`
}

// UserProfileResponse 用户资料响应结构
type UserProfileResponse struct {
	ID       int64  `json:"id"`
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Timezone string `json:"timezone"`
	Language string `json:"language"`
	Theme    string `json:"theme"`
	Role     string `json:"role"`
	Points   int    `json:"points"`
}

// GetSettings 获取用户设置
func (h *GinSettingsHandler) GetSettings(c *gin.Context) {
	// 返回默认设置，实际项目中应该从数据库获取用户设置
	settings := SettingsResponse{
		Theme: struct {
			Mode  string `json:"mode"`
			Color string `json:"color"`
		}{
			Mode:  "light",
			Color: "blue",
		},
		Language:   "zh-CN",
		DateFormat: "yyyy-MM-dd",
		TimeFormat: "HH:mm:ss",
		Timezone:   "Asia/Shanghai",
		Layout: struct {
			SidebarPosition string `json:"sidebarPosition"`
			ContentWidth    string `json:"contentWidth"`
			CompactMode     bool   `json:"compactMode"`
		}{
			SidebarPosition: "left",
			ContentWidth:    "contained",
			CompactMode:     false,
		},
		Notifications: struct {
			Email struct {
				Enabled   bool   `json:"enabled"`
				Frequency string `json:"frequency"`
				Types     struct {
					WorkflowCompleted bool `json:"workflow_completed"`
					WorkflowFailed    bool `json:"workflow_failed"`
					AgentError        bool `json:"agent_error"`
					NewUpdates        bool `json:"new_updates"`
					SecurityAlerts    bool `json:"security_alerts"`
				} `json:"types"`
			} `json:"email"`
			Browser struct {
				Enabled bool   `json:"enabled"`
				Sound   string `json:"sound"`
				Types   struct {
					WorkflowCompleted bool `json:"workflow_completed"`
					WorkflowFailed    bool `json:"workflow_failed"`
					AgentError        bool `json:"agent_error"`
					SecurityAlerts    bool `json:"security_alerts"`
				} `json:"types"`
			} `json:"browser"`
			InApp struct {
				Enabled  bool   `json:"enabled"`
				Position string `json:"position"`
				Types    struct {
					WorkflowCompleted bool `json:"workflow_completed"`
					WorkflowFailed    bool `json:"workflow_failed"`
					AgentError        bool `json:"agent_error"`
					SecurityAlerts    bool `json:"security_alerts"`
					NewUpdates        bool `json:"new_updates"`
				} `json:"types"`
			} `json:"inApp"`
		}{
			Email: struct {
				Enabled   bool   `json:"enabled"`
				Frequency string `json:"frequency"`
				Types     struct {
					WorkflowCompleted bool `json:"workflow_completed"`
					WorkflowFailed    bool `json:"workflow_failed"`
					AgentError        bool `json:"agent_error"`
					NewUpdates        bool `json:"new_updates"`
					SecurityAlerts    bool `json:"security_alerts"`
				} `json:"types"`
			}{
				Enabled:   true,
				Frequency: "immediate",
				Types: struct {
					WorkflowCompleted bool `json:"workflow_completed"`
					WorkflowFailed    bool `json:"workflow_failed"`
					AgentError        bool `json:"agent_error"`
					NewUpdates        bool `json:"new_updates"`
					SecurityAlerts    bool `json:"security_alerts"`
				}{
					WorkflowCompleted: true,
					WorkflowFailed:    true,
					AgentError:        true,
					NewUpdates:        true,
					SecurityAlerts:    true,
				},
			},
			Browser: struct {
				Enabled bool   `json:"enabled"`
				Sound   string `json:"sound"`
				Types   struct {
					WorkflowCompleted bool `json:"workflow_completed"`
					WorkflowFailed    bool `json:"workflow_failed"`
					AgentError        bool `json:"agent_error"`
					SecurityAlerts    bool `json:"security_alerts"`
				} `json:"types"`
			}{
				Enabled: false,
				Sound:   "default",
				Types: struct {
					WorkflowCompleted bool `json:"workflow_completed"`
					WorkflowFailed    bool `json:"workflow_failed"`
					AgentError        bool `json:"agent_error"`
					SecurityAlerts    bool `json:"security_alerts"`
				}{
					WorkflowCompleted: true,
					WorkflowFailed:    true,
					AgentError:        true,
					SecurityAlerts:    true,
				},
			},
			InApp: struct {
				Enabled  bool   `json:"enabled"`
				Position string `json:"position"`
				Types    struct {
					WorkflowCompleted bool `json:"workflow_completed"`
					WorkflowFailed    bool `json:"workflow_failed"`
					AgentError        bool `json:"agent_error"`
					SecurityAlerts    bool `json:"security_alerts"`
					NewUpdates        bool `json:"new_updates"`
				} `json:"types"`
			}{
				Enabled:  true,
				Position: "top-right",
				Types: struct {
					WorkflowCompleted bool `json:"workflow_completed"`
					WorkflowFailed    bool `json:"workflow_failed"`
					AgentError        bool `json:"agent_error"`
					SecurityAlerts    bool `json:"security_alerts"`
					NewUpdates        bool `json:"new_updates"`
				}{
					WorkflowCompleted: true,
					WorkflowFailed:    true,
					AgentError:        true,
					SecurityAlerts:    true,
					NewUpdates:        true,
				},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   settings,
	})
}

// UpdateSettings 更新用户设置
func (h *GinSettingsHandler) UpdateSettings(c *gin.Context) {
	var settings SettingsResponse
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 这里应该保存设置到数据库
	// 目前只是返回成功响应
	h.Logger.Info("用户设置已更新", zap.Any("settings", settings))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   settings,
	})
}

// GetUserProfile 获取用户资料
func (h *GinSettingsHandler) GetUserProfile(c *gin.Context) {
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
			return
		}
		h.Logger.Error("获取用户资料失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取用户资料失败: " + err.Error(),
		})
		return
	}

	// 构建响应数据
	profile := UserProfileResponse{
		ID:       user.ID,
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		Timezone: user.Timezone,
		Language: user.Language,
		Theme:    user.Theme,
		Role:     user.Role,
		Points:   user.Points,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   profile,
	})
}

// UpdateUserProfile 更新用户资料
func (h *GinSettingsHandler) UpdateUserProfile(c *gin.Context) {
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

	// 解析请求数据
	var input models.UserUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
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
			return
		}
		h.Logger.Error("获取用户信息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取用户信息失败: " + err.Error(),
		})
		return
	}

	// 更新用户资料
	if err := user.Update(h.DB, input); err != nil {
		if err == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "邮箱已被使用",
			})
			return
		}
		h.Logger.Error("更新用户资料失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新用户资料失败: " + err.Error(),
		})
		return
	}

	// 返回更新后的用户资料
	profile := UserProfileResponse{
		ID:       user.ID,
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		Timezone: user.Timezone,
		Language: user.Language,
		Theme:    user.Theme,
		Role:     user.Role,
		Points:   user.Points,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "用户资料更新成功",
		"data":    profile,
	})
}

// UploadAvatar 上传头像
func (h *GinSettingsHandler) UploadAvatar(c *gin.Context) {
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
	filename := "avatar_" + uid + "_" + strconv.FormatInt(c.Request.Context().Value("timestamp").(int64), 10) + ".jpg"

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

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "头像上传成功",
		"data": gin.H{
			"avatar": avatarURL,
		},
	})
}

// ChangePassword 修改密码
func (h *GinSettingsHandler) ChangePassword(c *gin.Context) {
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

	// 解析请求数据
	var input models.PasswordChangeInput
	if err := c.ShouldBindJSON(&input); err != nil {
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
			return
		}
		h.Logger.Error("获取用户信息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "获取用户信息失败: " + err.Error(),
		})
		return
	}

	// 修改密码
	if err := user.ChangePassword(h.DB, input); err != nil {
		if err == gorm.ErrInvalidValue {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "当前密码不正确",
			})
			return
		}
		h.Logger.Error("修改密码失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "修改密码失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "密码修改成功",
	})
}
