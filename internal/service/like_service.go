package service

import (
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/logger"

	"gorm.io/gorm"
)

type LikeService struct {
	likeRepo    *repository.LikeRepository
	articleRepo *repository.ArticleRepository
	userRepo    *repository.UserRepository
	db          *gorm.DB
}

func NewLikeService(
	likeRepo *repository.LikeRepository,
	articleRepo *repository.ArticleRepository,
	userRepo *repository.UserRepository,
	db *gorm.DB,
) *LikeService {
	return &LikeService{
		likeRepo:    likeRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
		db:          db,
	}
}

// LikeArticle 点赞文章
func (s *LikeService) LikeArticle(articleID, userID uint) error {
	// 1. 检查文章是否存在
	_, err := s.articleRepo.GetByID(articleID)
	if err != nil {
		logger.Warnf("Like article failed: article not found, articleID=%d", articleID)
		return errors.New(errors.ArticleNotFoundCode)
	}

	// 2. 检查用户是否存在
	_, err = s.userRepo.GetByID(userID)
	if err != nil {
		logger.Warnf("Like article failed: user not found, userID=%d", userID)
		return errors.New(errors.UserNotFoundCode)
	}

	// 3. 检查是否已点赞
	if s.likeRepo.CheckLiked(articleID, userID) {
		logger.Infof("Like article skipped: already liked, articleID=%d, userID=%d", articleID, userID)
		return errors.New(errors.AlreadyLikedCode)
	}

	// 4. 使用事务创建点赞记录并增加计数
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建点赞记录
		like := &entity.Like{
			ArticleID: articleID,
			UserID:    userID,
		}
		if err := tx.Create(like).Error; err != nil {
			return err
		}

		// 增加文章点赞数
		if err := tx.Model(&entity.Article{}).
			Where("id = ?", articleID).
			Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Like article failed: database error, articleID=%d, userID=%d, error=%v", articleID, userID, err)
		return err
	}

	return nil
}

// UnlikeArticle 取消点赞文章
func (s *LikeService) UnlikeArticle(articleID, userID uint) error {
	// 1. 检查是否已点赞
	if !s.likeRepo.CheckLiked(articleID, userID) {
		logger.Infof("Unlike article skipped: not liked yet, articleID=%d, userID=%d", articleID, userID)
		return errors.New(errors.NotLikedYetCode)
	}

	// 2. 使用事务删除点赞记录并减少计数
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 删除点赞记录
		if err := tx.Where("article_id = ? AND user_id = ?", articleID, userID).
			Delete(&entity.Like{}).Error; err != nil {
			return err
		}

		// 减少文章点赞数
		if err := tx.Model(&entity.Article{}).
			Where("id = ?", articleID).
			Update("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Unlike article failed: database error, articleID=%d, userID=%d, error=%v", articleID, userID, err)
		return err
	}

	return nil
}

// GetLikeStatus 获取点赞状态
func (s *LikeService) GetLikeStatus(articleID, userID uint) (bool, error) {
	return s.likeRepo.CheckLiked(articleID, userID), nil
}
