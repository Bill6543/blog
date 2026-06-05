package service

import (
	"blog/internal/model/dto"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/logger"

	"gorm.io/gorm"
)

// 评论相关错误定义已移至 pkg/errors，使用统一错误码

type CommentService struct {
	commentRepo *repository.CommentRepository
	articleRepo *repository.ArticleRepository
	db          *gorm.DB
}

func NewCommentService(
	commentRepo *repository.CommentRepository,
	articleRepo *repository.ArticleRepository,
	db *gorm.DB,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		db:          db,
	}
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(req dto.CreateCommentRequest, userID uint) (*dto.CommentResponse, error) {
	// 验证文章是否存在
	_, err := s.articleRepo.GetByID(req.ArticleID)
	if err != nil {
		logger.Warnf("Create comment failed: article not found, articleID=%d", req.ArticleID)
		return nil, errors.New(errors.ArticleNotFoundCode)
	}

	// 验证父评论（如果提供了）
	if req.ParentID != nil && *req.ParentID > 0 {
		_, err := s.commentRepo.GetByID(*req.ParentID)
		if err != nil {
			logger.Warnf("Create comment failed: parent comment not found, parentID=%d", *req.ParentID)
			return nil, errors.New(errors.CommentNotFoundCode)
		}
	}

	comment := &entity.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID,
		Content:   req.Content,
		ParentID:  req.ParentID,
		Status:    1, // 待审核
	}

	// 创建评论（审核通过时才增加计数）
	if err := s.commentRepo.Create(comment); err != nil {
		logger.Errorf("Create comment failed: database error, articleID=%d, userID=%d, error=%v", req.ArticleID, userID, err)
		return nil, err
	}

	return s.convertToDTO(comment), nil
}

// GetCommentsByArticleID 获取文章评论列表
func (s *CommentService) GetCommentsByArticleID(articleID uint) ([]dto.CommentResponse, error) {
	comments, err := s.commentRepo.GetByArticleID(articleID)
	if err != nil {
		return nil, err
	}

	return s.convertToDTOs(comments), nil
}

// ApproveComment 审核评论（通过）
func (s *CommentService) ApproveComment(id uint) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return errors.New(errors.CommentNotFoundCode)
	}

	// 使用事务更新评论状态并增加计数
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 更新评论状态
		if err := tx.Model(&entity.Comment{}).
			Where("id = ?", id).
			Update("status", 2).Error; err != nil {
			return err
		}

		// 增加文章评论数
		if err := tx.Model(&entity.Article{}).
			Where("id = ?", comment.ArticleID).
			Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Approve comment failed: database error, commentID=%d, error=%v", id, err)
		return err
	}

	return nil
}

// RejectComment 审核评论（拒绝）
func (s *CommentService) RejectComment(id uint) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return errors.New(errors.CommentNotFoundCode)
	}

	comment.Status = -1 // 拒绝
	return s.commentRepo.Update(comment)
}

// GetCommentByID 获取评论详情
func (s *CommentService) GetCommentByID(id uint) (*dto.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(errors.CommentNotFoundCode)
	}
	return s.convertToDTO(comment), nil
}

// GetPendingComments 获取待审核评论列表
func (s *CommentService) GetPendingComments() ([]dto.CommentResponse, error) {
	comments, err := s.commentRepo.GetPendingComments()
	if err != nil {
		return nil, err
	}
	return s.convertToDTOs(comments), nil
}

// GetPendingCommentCount 获取待审核评论数量
func (s *CommentService) GetPendingCommentCount() (int64, error) {
	var count int64
	err := s.db.Model(&entity.Comment{}).Where("status = ?", 1).Count(&count).Error
	return count, err
}

// UpdateComment 更新评论（支持管理员和普通用户）
func (s *CommentService) UpdateComment(id uint, req dto.UpdateCommentRequest, userID uint, isAdmin bool) (*dto.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(errors.CommentNotFoundCode)
	}

	// 验证权限：管理员或作者
	if !isAdmin && comment.UserID != userID {
		logger.Warnf("Update comment failed: permission denied, commentID=%d, authorID=%d, requestUserID=%d, isAdmin=%v",
			id, comment.UserID, userID, isAdmin)
		return nil, errors.New(errors.PermissionDenied)
	}

	// 更新内容
	if req.Content != "" {
		comment.Content = req.Content
	}

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, err
	}

	return s.convertToDTO(comment), nil
}

// DeleteComment 删除评论（支持管理员和普通用户）
func (s *CommentService) DeleteComment(id uint, userID uint, isAdmin bool) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return errors.New(errors.CommentNotFoundCode)
	}

	// 验证权限：管理员或作者
	if !isAdmin && comment.UserID != userID {
		logger.Warnf("Delete comment failed: permission denied, commentID=%d, authorID=%d, requestUserID=%d, isAdmin=%v",
			id, comment.UserID, userID, isAdmin)
		return errors.New(errors.PermissionDenied)
	}

	// 使用事务删除评论并减少计数
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 软删除评论
		if err := tx.Model(&entity.Comment{}).
			Where("id = ?", id).
			Update("status", -1).Error; err != nil {
			return err
		}

		// 减少文章评论数（使用 GREATEST 确保不会出现负数）
		if err := tx.Model(&entity.Article{}).
			Where("id = ?", comment.ArticleID).
			Update("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)")).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Delete comment failed: database error, commentID=%d, error=%v", id, err)
		return err
	}

	return nil
}

// convertToDTO 转换为 DTO（防御式处理关联数据）
func (s *CommentService) convertToDTO(comment *entity.Comment) *dto.CommentResponse {
	return s.convertToDTOWithReplies(comment, true)
}

// convertToDTOWithReplies 转换为 DTO，可选择是否加载 replies
func (s *CommentService) convertToDTOWithReplies(comment *entity.Comment, loadReplies bool) *dto.CommentResponse {
	commentDTO := &dto.CommentResponse{
		ID:        comment.ID,
		ArticleID: comment.ArticleID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		Status:    comment.Status,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	// 防御式处理：检查 Article 是否加载成功
	if comment.Article.ID > 0 {
		commentDTO.Article = &dto.ArticleResponse{
			ID:     comment.Article.ID,
			Title:  comment.Article.Title,
			Status: comment.Article.Status,
		}

		// 如果有分类信息，转换分类
		if comment.Article.Category.ID > 0 {
			commentDTO.Article.Category = &dto.CategoryDTO{
				ID:   comment.Article.Category.ID,
				Name: comment.Article.Category.Name,
			}
		}

		// 如果有标签信息，转换标签
		if len(comment.Article.Tags) > 0 {
			for _, tag := range comment.Article.Tags {
				commentDTO.Article.Tags = append(commentDTO.Article.Tags, dto.TagDTO{
					ID:   tag.ID,
					Name: tag.Name,
				})
			}
		}
	}

	// 防御式处理：检查 User 是否加载成功
	if comment.User.ID > 0 {
		commentDTO.User = &dto.UserResponse{
			ID:       comment.User.ID,
			Username: comment.User.Username,
			Nickname: comment.User.Nickname,
			Avatar:   comment.User.Avatar,
			Bio:      comment.User.Bio,
			Role:     comment.User.Role,
		}
	}

	// 防御式处理：检查 Parent 是否加载成功（不递归加载 Parent 的 replies）
	if comment.Parent != nil && comment.Parent.ID > 0 {
		commentDTO.Parent = s.convertToDTOWithReplies(comment.Parent, false)
	}

	// 如果是父评论且需要加载回复，递归加载回复
	if loadReplies && comment.ParentID == nil {
		replies, _ := s.commentRepo.GetReplies(comment.ID)
		for _, reply := range replies {
			commentDTO.Replies = append(commentDTO.Replies, *s.convertToDTOWithReplies(&reply, false))
		}
	}

	return commentDTO
}

// convertToDTOs 批量转换
func (s *CommentService) convertToDTOs(comments []entity.Comment) []dto.CommentResponse {
	var dtos []dto.CommentResponse
	for _, comment := range comments {
		dto := s.convertToDTO(&comment)
		if dto != nil {
			dtos = append(dtos, *dto)
		}
	}
	return dtos
}
