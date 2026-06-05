package dto

import "time"

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title      string `form:"title" binding:"required,min=1,max=200"`
	Summary    string `form:"summary" binding:"omitempty,max=500"`
	Content    string `form:"content" binding:"required,min=1,max=100000"`
	CoverImage string `form:"cover_image"`
	CategoryID uint   `form:"category_id" binding:"omitempty,gt=0"`
	TagIDs     []uint `form:"tag_ids"`
	Status     string `form:"status" binding:"omitempty,oneof=0 1"` // 0:草稿, 1:发布（使用字符串避免 FormData 中 0 被视为未传递）
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title      string `form:"title" binding:"omitempty,min=1,max=200"`
	Summary    string `form:"summary" binding:"omitempty,max=500"`
	Content    string `form:"content" binding:"omitempty,min=1,max=100000"`
	CoverImage string `form:"cover_image"`
	CategoryID uint   `form:"category_id" binding:"omitempty,gt=0"`
	TagIDs     []uint `form:"tag_ids"`
	Status     string `form:"status" binding:"omitempty,oneof=0 1"` // 0:草稿, 1:发布（使用字符串避免 FormData 中 0 被视为未传递）
}

// ArticleResponse 文章响应
type ArticleResponse struct {
	ID           uint          `json:"id"`
	Title        string        `json:"title"`
	Summary      string        `json:"summary"`
	Content      string        `json:"content"`
	CoverImage   string        `json:"cover_image"`
	AuthorID     uint          `json:"author_id"`
	Author       *UserResponse `json:"author"`
	CategoryID   uint          `json:"category_id"`
	Category     *CategoryDTO  `json:"category"`
	Tags         []TagDTO      `json:"tags"`
	ViewCount    int           `json:"view_count"`
	LikeCount    int           `json:"like_count"`
	CommentCount int           `json:"comment_count"`
	Status       int           `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// ArticleListResponse 文章列表响应
type ArticleListResponse struct {
	Total int64             `json:"total"`
	List  []ArticleResponse `json:"list"`
}
