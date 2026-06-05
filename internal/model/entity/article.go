package entity

import (
	"time"
)

// Article 文章实体
type Article struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	Title        string    `gorm:"size:200;not null" json:"title"`
	Summary      string    `gorm:"type:text" json:"summary"`
	Content      string    `gorm:"type:longtext;not null" json:"content"`
	CoverImage   string    `gorm:"size:255" json:"cover_image"`
	AuthorID     uint      `gorm:"not null" json:"author_id"`
	Author       User      `gorm:"foreignKey:AuthorID" json:"author"`
	CategoryID   uint      `json:"category_id"`
	Category     Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Tags         []Tag     `gorm:"many2many:article_tag;" json:"tags"`
	ViewCount    int       `gorm:"default:0" json:"view_count"`
	LikeCount    int       `gorm:"default:0" json:"like_count"`
	CommentCount int       `gorm:"default:0" json:"comment_count"`
	Status       int       `json:"status"` // 1:发布，0:草稿
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}
