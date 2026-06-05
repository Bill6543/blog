package entity

import (
	"time"
)

// Tag 标签实体
type Tag struct {
	ID        uint      `gorm:"primary_key" json:"Id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Color     string    `gorm:"size:20" json:"color"`
	Status    int       `gorm:"default:1" json:"status"` // 1:启用，0:禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
