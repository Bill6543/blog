package entity

import (
	"time"
)

// Category 分类实体
type Category struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      int       `gorm:"default:1" json:"status"` // 1:启用，0:禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}
