package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // 不暴露密码哈希
	FullName     string     `json:"fullName"`
	Avatar       string     `json:"avatar"`
	Role         string     `json:"role"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`
}

// UserRegisterInput 用户注册输入
type UserRegisterInput struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"fullName" validate:"required"`
}

// UserLoginInput 用户登录输入
type UserLoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserUpdateInput 用户更新输入
type UserUpdateInput struct {
	Email    string `json:"email" validate:"omitempty,email"`
	FullName string `json:"fullName"`
	Avatar   string `json:"avatar"`
}

// PasswordChangeInput 密码更改输入
type PasswordChangeInput struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
}

// HashPassword 创建密码哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// CreateUser 创建新用户
func CreateUser(db *sql.DB, input UserRegisterInput) (*User, error) {
	// TODO: 实现用户创建逻辑
	return nil, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	// TODO: 实现根据ID获取用户
	return nil, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	// TODO: 实现根据用户名获取用户
	return nil, nil
}

// UpdateUser 更新用户信息
func (u *User) Update(db *sql.DB, input UserUpdateInput) error {
	// TODO: 实现用户更新逻辑
	return nil
}

// ChangePassword 修改用户密码
func (u *User) ChangePassword(db *sql.DB, input PasswordChangeInput) error {
	// TODO: 实现密码更改逻辑
	return nil
}
