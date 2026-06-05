package entity

import (
	"time"
)

// User 用户实体
type User struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	Username     string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password     string    `gorm:"size:255;not null" json:"-"`
	Nickname     string    `gorm:"size:50" json:"nickname"`
	Avatar       string    `gorm:"size:255" json:"avatar"`             // 头像
	Bio          string    `gorm:"type:text" json:"bio"`               // 个人简介
	Role         string    `gorm:"size:20;default:'user'" json:"role"` // admin, user
	Status       int       `gorm:"default:1" json:"status"`            // 1:正常，0:禁用
	LoginSession string    `gorm:"size:255" json:"login_session"`      //当前有效的登录 session
	LastLoginTime *time.Time `gorm:"index" json:"last_login_time"`     //最后登录时间（用于并发控制）
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
