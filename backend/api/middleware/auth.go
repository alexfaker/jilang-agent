package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexfaker/jilang-agent/utils"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware 认证中间件，验证JWT令牌
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从请求头获取令牌
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "缺少认证令牌")
				return
			}

			// 提取令牌
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
				utils.RespondWithError(w, http.StatusUnauthorized, "认证格式无效")
				return
			}

			tokenString := bearerToken[1]

			// 解析令牌
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// 验证签名算法
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					utils.RespondWithError(w, http.StatusUnauthorized, "无效的令牌签名")
					return
				}
				utils.RespondWithError(w, http.StatusUnauthorized, "无法验证令牌: "+err.Error())
				return
			}

			if !token.Valid {
				utils.RespondWithError(w, http.StatusUnauthorized, "无效的令牌")
				return
			}

			// 从令牌中提取声明
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "无效的令牌声明")
				return
			}

			// 提取用户信息
			userID, ok := claims["user_id"].(float64)
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户标识")
				return
			}

			username, ok := claims["username"].(string)
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户名")
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "无效的用户角色")
				return
			}

			// 将用户信息添加到请求上下文
			ctx := context.WithValue(r.Context(), "userID", int64(userID))
			ctx = context.WithValue(ctx, "username", username)
			ctx = context.WithValue(ctx, "role", role)

			// 使用更新后的上下文继续处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
