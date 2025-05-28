package handlers

import (
	"net/http"

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
