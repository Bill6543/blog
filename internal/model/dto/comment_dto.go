package dto

import "time"

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID        uint              `json:"id"`
	ArticleID uint              `json:"article_id"`
	Article   *ArticleResponse  `json:"article"`
	UserID    uint              `json:"user_id"`
	User      *UserResponse     `json:"user"`
	ParentID  *uint             `json:"parent_id"`
	Parent    *CommentResponse  `json:"parent"`
	Content   string            `json:"content"`
	Status    int               `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Replies   []CommentResponse `json:"replies"`
}
