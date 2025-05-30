package models

import (
	"time"

	"github.com/alexfaker/jilang-agent/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       string     `json:"userID" gorm:"type:varchar(50);uniqueIndex;not null"`
	Username     string     `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email        string     `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string     `json:"-" gorm:"column:password_hash;type:varchar(255);not null"` // 不暴露密码哈希
	FullName     string     `json:"fullName" gorm:"column:full_name;type:varchar(100)"`
	Avatar       string     `json:"avatar" gorm:"type:varchar(255)"`
	Bio          string     `json:"bio" gorm:"type:text"`                                     // 个人简介
	Timezone     string     `json:"timezone" gorm:"type:varchar(50);default:'Asia/Shanghai'"` // 时区
	Language     string     `json:"language" gorm:"type:varchar(10);default:'zh_CN'"`         // 语言
	Theme        string     `json:"theme" gorm:"type:varchar(20);default:'light'"`            // 主题
	Role         string     `json:"role" gorm:"type:varchar(20);default:'user'"`
	Points       int        `json:"points" gorm:"default:0;not null"` // 用户点数余额
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	LastLoginAt  *time.Time `json:"lastLoginAt" gorm:"column:last_login_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
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
	Bio      string `json:"bio"`
	Timezone string `json:"timezone"`
	Language string `json:"language"`
	Theme    string `json:"theme"`
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

// CreateUser 使用GORM创建新用户
func CreateUser(db *gorm.DB, input UserRegisterInput) (*User, error) {
	// 检查用户名是否已存在
	var count int64
	if err := db.Model(&User{}).Where("username = ?", input.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, gorm.ErrDuplicatedKey
	}

	// 检查邮箱是否已存在
	if err := db.Model(&User{}).Where("email = ?", input.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, gorm.ErrDuplicatedKey
	}

	// 生成密码哈希
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// 生成全局唯一的用户ID，最多重试3次
	var userID string
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		userID = utils.GenerateUserID()

		// 检查UserID是否已存在
		var existingUser User
		err := db.Where("user_id = ?", userID).First(&existingUser).Error
		if err == gorm.ErrRecordNotFound {
			// UserID不存在，可以使用
			break
		} else if err != nil {
			// 查询出错
			return nil, err
		}

		// 如果是最后一次重试仍然重复，返回错误
		if i == maxRetries-1 {
			return nil, gorm.ErrDuplicatedKey
		}
	}

	// 创建用户对象
	user := &User{
		UserID:       userID,
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: passwordHash,
		FullName:     input.FullName,
		Avatar:       "/static/avatars/default.png",
		Role:         "user",
		Points:       0,
	}

	// 保存到数据库
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID 使用GORM根据ID获取用户
func GetUserByID(db *gorm.DB, id int64) (*User, error) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 使用GORM根据用户名获取用户
func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 使用GORM更新用户信息
func (u *User) Update(db *gorm.DB, input UserUpdateInput) error {
	// 检查邮箱是否已被其他用户使用
	if input.Email != "" && input.Email != u.Email {
		var count int64
		if err := db.Model(&User{}).Where("email = ? AND id != ?", input.Email, u.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return gorm.ErrDuplicatedKey
		}
	}

	// 准备更新数据
	updates := map[string]interface{}{}
	if input.Email != "" {
		updates["email"] = input.Email
		u.Email = input.Email
	}
	if input.FullName != "" {
		updates["full_name"] = input.FullName
		u.FullName = input.FullName
	}
	if input.Avatar != "" {
		updates["avatar"] = input.Avatar
		u.Avatar = input.Avatar
	}
	if input.Bio != "" {
		updates["bio"] = input.Bio
		u.Bio = input.Bio
	}
	if input.Timezone != "" {
		updates["timezone"] = input.Timezone
		u.Timezone = input.Timezone
	}
	if input.Language != "" {
		updates["language"] = input.Language
		u.Language = input.Language
	}
	if input.Theme != "" {
		updates["theme"] = input.Theme
		u.Theme = input.Theme
	}

	// 更新用户信息
	return db.Model(u).Updates(updates).Error
}

// ChangePassword 使用GORM修改用户密码
func (u *User) ChangePassword(db *gorm.DB, input PasswordChangeInput) error {
	// 验证当前密码
	if !u.CheckPassword(input.CurrentPassword) {
		return gorm.ErrInvalidValue
	}

	// 生成新密码哈希
	newPasswordHash, err := HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	// 更新密码
	if err := db.Model(u).Updates(map[string]interface{}{
		"password_hash": newPasswordHash,
	}).Error; err != nil {
		return err
	}

	// 更新内存中的密码哈希
	u.PasswordHash = newPasswordHash

	return nil
}

// UpdateLastLogin 使用GORM更新最后登录时间
func (u *User) UpdateLastLogin(db *gorm.DB) error {
	now := time.Now()
	if err := db.Model(u).Update("last_login_at", now).Error; err != nil {
		return err
	}

	u.LastLoginAt = &now
	return nil
}
