package entity

import (
	"time"
)

// Like 点赞实体
type Like struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	ArticleID uint      `gorm:"not null;uniqueIndex:uk_like_article_user" json:"article_id"`
	UserID    uint      `gorm:"not null;uniqueIndex:uk_like_article_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}
