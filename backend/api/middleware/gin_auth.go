package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// RequestIDMiddleware 为请求添加请求ID
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取请求ID，如果没有则生成一个新的
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String() // 使用 UUID 替代 utils.GenerateRequestID
		}

		// 设置请求ID到上下文和响应头
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// GinAuthMiddleware 验证JWT令牌的中间件
func GinAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "未提供授权令牌",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "授权格式无效",
			})
			c.Abort()
			return
		}

		// 获取令牌
		tokenString := parts[1]

		// 解析令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		// 处理解析错误
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "无效的令牌: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 验证令牌有效性
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "无效的令牌",
			})
			c.Abort()
			return
		}

		// 从令牌中获取声明
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "无效的令牌声明",
			})
			c.Abort()
			return
		}

		// 获取用户ID和用户名
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "无效的用户ID",
			})
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "无效的用户名",
			})
			c.Abort()
			return
		}

		// 将用户信息设置到上下文
		c.Set("userID", userID)
		c.Set("username", username)

		c.Next()
	}
}
