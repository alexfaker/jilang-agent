package handlers

import (
	"net/http"
	"time"

	"github.com/alexfaker/jilang-agent/config"
	"github.com/alexfaker/jilang-agent/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GinAuthHandler 处理认证相关的请求
type GinAuthHandler struct {
	DB        *gorm.DB
	Logger    *zap.Logger
	Config    config.AuthConfig
	Validator *validator.Validate
}

// NewGinAuthHandler 创建一个新的GinAuthHandler实例
func NewGinAuthHandler(db *gorm.DB, logger *zap.Logger, cfg config.AuthConfig) *GinAuthHandler {
	return &GinAuthHandler{
		DB:        db,
		Logger:    logger,
		Config:    cfg,
		Validator: validator.New(),
	}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register 用户注册
func (h *GinAuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := h.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "用户名已存在",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		h.Logger.Error("检查用户名失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "注册失败: " + err.Error(),
		})
		return
	}

	// 检查邮箱是否已存在
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "邮箱已被使用",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		h.Logger.Error("检查邮箱失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "注册失败: " + err.Error(),
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("密码加密失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "注册失败: " + err.Error(),
		})
		return
	}

	// 创建用户
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user", // 默认角色
	}

	if err := h.DB.Create(&user).Error; err != nil {
		h.Logger.Error("创建用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "注册失败: " + err.Error(),
		})
		return
	}

	// 生成令牌
	token, err := h.generateToken(user)
	if err != nil {
		h.Logger.Error("生成令牌失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "注册成功，但生成令牌失败: " + err.Error(),
		})
		return
	}

	// 返回用户信息和令牌
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
			"token": token,
		},
	})
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (h *GinAuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 查找用户 - 支持邮箱登录
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "邮箱或密码不正确",
			})
		} else {
			h.Logger.Error("查找用户失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "登录失败: " + err.Error(),
			})
		}
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "邮箱或密码不正确",
		})
		return
	}

	// 生成令牌
	token, err := h.generateToken(user)
	if err != nil {
		h.Logger.Error("生成令牌失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "登录失败: " + err.Error(),
		})
		return
	}

	// 更新最后登录时间
	if err := h.DB.Model(&user).Update("last_login", time.Now()).Error; err != nil {
		h.Logger.Warn("更新最后登录时间失败", zap.Error(err))
	}

	// 返回用户信息和令牌
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
			"token": token,
		},
	})
}

// RefreshTokenRequest 刷新令牌请求结构
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// RefreshToken 刷新令牌
func (h *GinAuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 解析令牌
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWTSecret), nil
	})
	// 处理解析错误
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的令牌: " + err.Error(),
		})
		return
	}

	// 验证令牌
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的令牌",
		})
		return
	}

	// 从令牌中获取声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的令牌声明",
		})
		return
	}

	// 获取用户ID
	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "无效的用户ID",
		})
		return
	}

	// 查找用户
	var user models.User
	if err := h.DB.First(&user, int64(userID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "用户不存在",
			})
		} else {
			h.Logger.Error("查找用户失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "刷新令牌失败: " + err.Error(),
			})
		}
		return
	}

	// 生成新令牌
	newToken, err := h.generateToken(user)
	if err != nil {
		h.Logger.Error("生成新令牌失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "刷新令牌失败: " + err.Error(),
		})
		return
	}

	// 返回新令牌
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token": newToken,
		},
	})
}

// generateToken 生成JWT令牌
func (h *GinAuthHandler) generateToken(user models.User) (string, error) {
	// 设置令牌过期时间
	expirationTime := time.Now().Add(time.Duration(h.Config.TokenExpiration) * time.Hour)

	// 创建令牌声明
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
