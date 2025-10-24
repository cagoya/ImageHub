package model

import (
	"time"
)

// User 用户信息
type User struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`                    // 主键，用户 ID
	Username     string    `json:"username" gorm:"uniqueIndex;type:varchar(64);not null"` // 用户名
	PasswordHash string    `json:"-" gorm:"type:varchar(64);not null"`                    // 密码哈希
	Email        string    `json:"email" gorm:"type:varchar(64);uniqueIndex;not null"`    // 邮箱
	IsAdmin      bool      `json:"isAdmin" gorm:"default:false"`                          // 是否为管理员
	CreatedTime  time.Time `json:"created_time" gorm:"autoCreateTime"`                    // 创建时间
	UpdatedTime  time.Time `json:"updated_time" gorm:"autoUpdateTime"`                    // 更新时间
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ForgetPasswordRequest 忘记密码请求
type ForgetPasswordRequest struct {
	NewPassword string `json:"new_password"`
	Username    string `json:"username"`
	Email       string `json:"email"`
}
