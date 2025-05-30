package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateUserID 生成全局唯一的用户ID
// 格式: USER_ + 无连字符的UUID
// 示例: USER_550e8400e29b11d4a716446655440000
func GenerateUserID() string {
	// 生成UUID并移除连字符
	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	return fmt.Sprintf("USER_%s", strings.ToUpper(uuidStr))
}

// GenerateShortUserID 生成较短的用户ID（仅用于展示）
// 格式: U + 8位随机字符
// 注意：此函数生成的ID较短但唯一性不如完整UUID
func GenerateShortUserID() string {
	// 取UUID的前8个字符（移除连字符后）
	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	return fmt.Sprintf("U%s", strings.ToUpper(uuidStr[:8]))
}

// ValidateUserID 验证用户ID格式是否正确
func ValidateUserID(userID string) bool {
	// 检查是否以USER_开头且长度正确
	if !strings.HasPrefix(userID, "USER_") {
		return false
	}

	// 移除前缀后应该是32位十六进制字符
	idPart := strings.TrimPrefix(userID, "USER_")
	if len(idPart) != 32 {
		return false
	}

	// 检查是否都是有效的十六进制字符
	for _, char := range idPart {
		if !((char >= '0' && char <= '9') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}

	return true
}
