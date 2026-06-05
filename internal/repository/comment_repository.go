package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *entity.Comment) error {
	return r.db.Create(comment).Error
}

// GetByID 根据 ID 查询评论
func (r *CommentRepository) GetByID(id uint) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.Preload("User").
		Preload("Parent").
		First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByArticleID 根据文章 ID 查询评论（只返回父评论）
func (r *CommentRepository) GetByArticleID(articleID uint) ([]entity.Comment, error) {
	var comments []entity.Comment
	err := r.db.Where("article_id = ? AND status = ? AND parent_id IS NULL", articleID, 2).
		Preload("User").
		Preload("Parent").
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

// Update 更新评论
func (r *CommentRepository) Update(comment *entity.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Model(&entity.Comment{}).
		Where("id = ?", id).
		Update("status", -1).Error
}

// GetReplies 获取评论回复
func (r *CommentRepository) GetReplies(parentID uint) ([]entity.Comment, error) {
	var replies []entity.Comment
	err := r.db.Where("parent_id = ? AND status = ?", parentID, 2).
		Preload("User").
		Preload("Parent").
		Preload("Parent.User").
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

// GetPendingComments 获取待审核评论列表
func (r *CommentRepository) GetPendingComments() ([]entity.Comment, error) {
	var comments []entity.Comment
	err := r.db.Where("status = ?", 1).
		Preload("User").
		Preload("Article").
		Preload("Parent").
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}
