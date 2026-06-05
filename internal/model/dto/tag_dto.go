package dto

import "time"

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=50"`
	Color string `json:"color"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name  string `json:"name" binding:"min=1,max=50"`
	Color string `json:"color"`
}

// TagDTO 标签 DTO
type TagDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
