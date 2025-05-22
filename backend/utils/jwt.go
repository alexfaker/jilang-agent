package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT 生成JWT令牌
func GenerateJWT(userID int64, username, role, jwtSecret string, expirationHours int) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// 创建令牌声明
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
