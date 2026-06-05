package repository

import (
	"blog/internal/model/entity"
	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

// Create 创建点赞记录
func (r *LikeRepository) Create(like *entity.Like) error {
	return r.db.Create(like).Error
}

// CheckLiked 检查是否已点赞
func (r *LikeRepository) CheckLiked(articleID, userID uint) bool {
	var count int64
	r.db.Model(&entity.Like{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Count(&count)
	return count > 0
}

// Delete 取消点赞
func (r *LikeRepository) Delete(articleID, userID uint) error {
	return r.db.Where("article_id = ? AND user_id = ?", articleID, userID).
		Delete(&entity.Like{}).Error
}

// GetCountByArticleID 获取文章点赞数
func (r *LikeRepository) GetCountByArticleID(articleID uint) int64 {
	var count int64
	r.db.Model(&entity.Like{}).
		Where("article_id = ?", articleID).
		Count(&count)
	return count
}
