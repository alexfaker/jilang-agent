package models

import (
	"database/sql"
	"errors"
	"fmt"
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
	// 检查用户名是否已存在
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", input.Username).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", input.Email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return nil, errors.New("邮箱已存在")
	}

	// 生成密码哈希
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希生成失败: %w", err)
	}

	// 默认头像URL
	defaultAvatar := "/static/avatars/default.png"

	// 插入新用户
	query := `
		INSERT INTO users (username, email, password_hash, full_name, avatar, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, 'user', NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	var user User
	err = db.QueryRow(
		query,
		input.Username,
		input.Email,
		passwordHash,
		input.FullName,
		defaultAvatar,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 设置返回对象的其他字段
	user.Username = input.Username
	user.Email = input.Email
	user.PasswordHash = passwordHash
	user.FullName = input.FullName
	user.Avatar = defaultAvatar
	user.Role = "user"

	return &user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, avatar, role, created_at, updated_at, last_login_at
		FROM users
		WHERE id = ?
	`
	var user User
	var lastLoginAt sql.NullTime

	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Avatar,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, avatar, role, created_at, updated_at, last_login_at
		FROM users
		WHERE username = ?
	`
	var user User
	var lastLoginAt sql.NullTime

	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Avatar,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (u *User) Update(db *sql.DB, input UserUpdateInput) error {
	// 检查邮箱是否已被其他用户使用
	if input.Email != "" && input.Email != u.Email {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)", input.Email, u.ID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("检查邮箱失败: %w", err)
		}
		if exists {
			return errors.New("邮箱已被其他用户使用")
		}
	}

	// 更新用户信息
	query := `
		UPDATE users
		SET 
			email = COALESCE(?, email),
			full_name = COALESCE(?, full_name),
			avatar = COALESCE(?, avatar),
			updated_at = NOW()
		WHERE id = ?
		RETURNING updated_at
	`

	err := db.QueryRow(
		query,
		nullString(input.Email),
		nullString(input.FullName),
		nullString(input.Avatar),
		u.ID,
	).Scan(&u.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("更新用户失败: %w", err)
	}

	// 更新内存中的用户对象
	if input.Email != "" {
		u.Email = input.Email
	}
	if input.FullName != "" {
		u.FullName = input.FullName
	}
	if input.Avatar != "" {
		u.Avatar = input.Avatar
	}

	return nil
}

// ChangePassword 修改用户密码
func (u *User) ChangePassword(db *sql.DB, input PasswordChangeInput) error {
	// 验证当前密码
	if !u.CheckPassword(input.CurrentPassword) {
		return errors.New("当前密码不正确")
	}

	// 生成新密码哈希
	newPasswordHash, err := HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("密码哈希生成失败: %w", err)
	}

	// 更新密码
	query := `
		UPDATE users
		SET 
			password_hash = ?,
			updated_at = NOW()
		WHERE id = ?
		RETURNING updated_at
	`

	err = db.QueryRow(query, newPasswordHash, u.ID).Scan(&u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 更新内存中的密码哈希
	u.PasswordHash = newPasswordHash

	return nil
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin(db *sql.DB) error {
	now := time.Now()
	query := `
		UPDATE users
		SET last_login_at = ?
		WHERE id = ?
	`

	_, err := db.Exec(query, now, u.ID)
	if err != nil {
		return fmt.Errorf("更新最后登录时间失败: %w", err)
	}

	u.LastLoginAt = &now
	return nil
}

// 辅助函数：将空字符串转为SQL NULL
func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
