package entity

import (
	"time"
)

// Comment 评论实体
type Comment struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	ArticleID uint      `gorm:"not null" json:"article_id"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"article"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ParentID  *uint     `json:"parent_id"`
	Parent    *Comment  `gorm:"foreignKey:ParentID" json:"parent"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Status    int       `gorm:"default:1" json:"status"` // 1:待审核，2:通过，-1:拒绝
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
