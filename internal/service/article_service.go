package service

import (
	"blog/internal/cache"
	"blog/internal/model/dto"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/logger"

	"gorm.io/gorm"
)

type ArticleService struct {
	articleRepo  *repository.ArticleRepository
	tagRepo      *repository.TagRepository
	categoryRepo *repository.CategoryRepository
	viewCounter  *cache.ViewCounter
	db           *gorm.DB
}

func NewArticleService(
	articleRepo *repository.ArticleRepository,
	tagRepo *repository.TagRepository,
	categoryRepo *repository.CategoryRepository,
	viewCounter *cache.ViewCounter,
	db *gorm.DB,
) *ArticleService {
	return &ArticleService{
		articleRepo:  articleRepo,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
		viewCounter:  viewCounter,
		db:           db,
	}
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(req dto.CreateArticleRequest, authorID uint) (*dto.ArticleResponse, error) {
	if req.CategoryID > 0 {
		_, err := s.categoryRepo.GetByID(req.CategoryID)
		if err != nil {
			logger.Warnf("Create article failed: category not found, categoryID=%d", req.CategoryID)
			return nil, errors.New(errors.CategoryNotFoundCode)
		}
	}

	// 2. 验证标签是否存在（如果提供了）
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.GetByIDs(req.TagIDs)
		if err != nil {
			logger.Warnf("Create article failed: tag not found, tagIDs=%v", req.TagIDs)
			return nil, errors.New(errors.TagNotFoundCode)
		}
		if len(tags) != len(req.TagIDs) {
			logger.Warnf("Create article failed: some tags not found, requested=%d, found=%d", len(req.TagIDs), len(tags))
			return nil, errors.New(errors.TagNotFoundCode)
		}
	}

	// 3. 创建文章
	coverImage := req.CoverImage
	if coverImage == "" {
		coverImage = "/static/default_cover.png" // 设置默认封面
	}

	// 解析状态（字符串转int）
	status := 1 // 默认发布状态
	if req.Status == "0" {
		status = 0 // 草稿
	} else if req.Status == "1" {
		status = 1 // 发布
	}

	article := &entity.Article{
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: coverImage,
		AuthorID:   authorID,
		CategoryID: req.CategoryID,
		Status:     status, // 使用解析后的状态（0:草稿, 1:发布）
	}

	// 4. 使用事务创建文章和关联标签
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 创建文章
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		// 关联标签
		if len(req.TagIDs) > 0 {
			for _, tagID := range req.TagIDs {
				if err := tx.Exec("INSERT INTO article_tag (article_id, tag_id) VALUES (?, ?)", article.ID, tagID).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Create article failed: database error, authorID=%d, title=%s, error=%v", authorID, req.Title, err)
		return nil, err
	}

	logger.Debugf("Article created successfully: articleID=%d, title=%s, authorID=%d", article.ID, req.Title, authorID)

	// 重新查询带关联数据的文章
	articleWithRelations, err := s.articleRepo.GetByID(article.ID)
	if err != nil {
		logger.Errorf("Failed to reload article with relations: articleID=%d, error=%v", article.ID, err)
		return s.convertToResponse(article), nil // 返回基本数据
	}

	return s.convertToResponse(articleWithRelations), nil
}

// GetArticleByID 获取文章详情
func (s *ArticleService) GetArticleByID(id uint) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		logger.Warnf("Get article failed: article not found, articleID=%d", id)
		return nil, errors.New(errors.ArticleNotFoundCode)
	}

	// 移除自动增加浏览量逻辑，改为由前端主动调用增加浏览量接口
	logger.Debugf("Article retrieved: articleID=%d, viewCount=%d", id, article.ViewCount)
	return s.convertToResponse(article), nil
}

// GetArticleList 获取文章列表
func (s *ArticleService) GetArticleList(page, pageSize int, status *int, categoryID *uint, tagID *uint) (*dto.ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var articles []entity.Article
	var total int64
	var err error

	// Use different repository methods based on filtering parameters
	if tagID != nil {
		articles, total, err = s.articleRepo.GetByTagID(*tagID, page, pageSize)
	} else if categoryID != nil {
		articles, total, err = s.articleRepo.GetByCategoryID(*categoryID, page, pageSize)
	} else {
		articles, total, err = s.articleRepo.GetList(page, pageSize, status)
	}

	if err != nil {
		return nil, err
	}

	// 转换为响应 DTO
	var list []dto.ArticleResponse
	for _, article := range articles {
		resp := s.convertToResponse(&article)
		list = append(list, *resp)
	}

	return &dto.ArticleListResponse{
		Total: total,
		List:  list,
	}, nil
}

// UpdateArticle 更新文章（支持管理员和普通用户）
func (s *ArticleService) UpdateArticle(id uint, req dto.UpdateArticleRequest, userID uint, isAdmin bool) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		logger.Warnf("Update article failed: article not found, articleID=%d", id)
		return nil, errors.New(errors.ArticleNotFoundCode)
	}

	// 验证权限：管理员或作者
	if !isAdmin && article.AuthorID != userID {
		logger.Warnf("Update article failed: permission denied, articleID=%d, authorID=%d, requestUserID=%d, isAdmin=%v",
			id, article.AuthorID, userID, isAdmin)
		return nil, errors.New(errors.PermissionDenied)
	}

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.CoverImage != "" {
		article.CoverImage = req.CoverImage
	}
	if req.CategoryID > 0 {
		article.CategoryID = req.CategoryID
	}
	// 支持更新状态（字符串转int）
	if req.Status == "0" {
		article.Status = 0 // 草稿
	} else if req.Status == "1" {
		article.Status = 1 // 发布
	}

	// 使用事务更新
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 更新文章
		if err := tx.Save(article).Error; err != nil {
			return err
		}

		// 更新标签关联
		if len(req.TagIDs) > 0 {
			// 删除旧关联
			if err := tx.Exec("DELETE FROM article_tag WHERE article_id = ?", id).Error; err != nil {
				return err
			}

			// 创建新关联
			for _, tagID := range req.TagIDs {
				if err := tx.Exec("INSERT INTO article_tag (article_id, tag_id) VALUES (?, ?)", id, tagID).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Update article failed: database error, articleID=%d, error=%v", id, err)
		return nil, err
	}

	logger.Debugf("Article updated successfully: articleID=%d, title=%s", id, article.Title)
	return s.convertToResponse(article), nil
}

// DeleteArticle 删除文章（支持管理员和普通用户）
func (s *ArticleService) DeleteArticle(id uint, userID uint, isAdmin bool) error {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		logger.Warnf("Delete article failed: article not found, articleID=%d", id)
		return errors.New(errors.ArticleNotFoundCode)
	}

	// 验证权限：管理员或作者
	if !isAdmin && article.AuthorID != userID {
		logger.Warnf("Delete article failed: permission denied, articleID=%d, authorID=%d, requestUserID=%d, isAdmin=%v",
			id, article.AuthorID, userID, isAdmin)
		return errors.New(errors.PermissionDenied)
	}

	if err := s.articleRepo.Delete(id); err != nil {
		logger.Errorf("Delete article failed: database error, articleID=%d, error=%v", id, err)
		return err
	}

	logger.Debugf("Article deleted successfully: articleID=%d, title=%s", id, article.Title)
	return nil
}

// RestoreArticle 恢复已删除的文章
func (s *ArticleService) RestoreArticle(id uint, userID uint, isAdmin bool) error {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		logger.Warnf("Restore article failed: article not found, articleID=%d", id)
		return errors.New(errors.ArticleNotFoundCode)
	}

	// 验证权限：管理员或作者
	if !isAdmin && article.AuthorID != userID {
		logger.Warnf("Restore article failed: permission denied, articleID=%d, authorID=%d, requestUserID=%d, isAdmin=%v",
			id, article.AuthorID, userID, isAdmin)
		return errors.New(errors.PermissionDenied)
	}

	// 检查文章是否为已删除状态
	if article.Status != -1 {
		logger.Warnf("Restore article failed: article is not deleted, articleID=%d, status=%d", id, article.Status)
		return errors.New(errors.ArticleNotFoundCode)
	}

	// 恢复文章（将状态改为草稿或发布，这里改为草稿）
	article.Status = 0
	if err := s.articleRepo.Update(article); err != nil {
		logger.Errorf("Restore article failed: database error, articleID=%d, error=%v", id, err)
		return err
	}

	logger.Debugf("Article restored successfully: articleID=%d, title=%s", id, article.Title)
	return nil
}

// IncreaseViewCount 增加文章浏览量（先累加到 Redis，定时批量落库）
func (s *ArticleService) IncreaseViewCount(id uint) error {
	if err := s.viewCounter.Incr(id); err != nil {
		// Redis 不可用时回退到直接写库，保证浏览量不丢
		logger.Warnf("View counter incr failed, fallback to DB: articleID=%d, error=%v", id, err)
		return s.articleRepo.IncrementViewCount(id)
	}
	logger.Debugf("View count increased: articleID=%d", id)
	return nil
}

// FlushViewCounts 将 Redis 中的浏览量增量批量落库（由后台定时任务调用）
func (s *ArticleService) FlushViewCounts() {
	counts, err := s.viewCounter.Drain()
	if err != nil {
		logger.Errorf("Flush view counts failed: drain error=%v", err)
		return
	}

	for articleID, n := range counts {
		if err := s.articleRepo.AddViewCount(articleID, n); err != nil {
			logger.Errorf("Flush view counts failed: articleID=%d, count=%d, error=%v", articleID, n, err)
		}
	}
}

// GetUserArticles 获取指定用户的文章列表
func (s *ArticleService) GetUserArticles(userID uint, page, pageSize int, status *int) (*dto.ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 调用 repository 的 GetByAuthorID 方法
	articles, total, err := s.articleRepo.GetByAuthorID(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 如果指定了状态筛选，需要在内存中过滤（因为 repository 层没有支持状态筛选）
	var filteredArticles []entity.Article
	if status != nil {
		for _, article := range articles {
			if article.Status == *status {
				filteredArticles = append(filteredArticles, article)
			}
		}
		articles = filteredArticles
	}

	// 转换为响应 DTO
	var list []dto.ArticleResponse
	for _, article := range articles {
		resp := s.convertToResponse(&article)
		list = append(list, *resp)
	}

	return &dto.ArticleListResponse{
		Total: total,
		List:  list,
	}, nil
}

// convertToResponse 转换为响应 DTO（防御式处理关联数据）
func (s *ArticleService) convertToResponse(article *entity.Article) *dto.ArticleResponse {
	resp := &dto.ArticleResponse{
		ID:           article.ID,
		Title:        article.Title,
		Summary:      article.Summary,
		Content:      article.Content,
		CoverImage:   article.CoverImage,
		AuthorID:     article.AuthorID,
		CategoryID:   article.CategoryID,
		ViewCount:    article.ViewCount,
		LikeCount:    article.LikeCount,
		CommentCount: article.CommentCount,
		Status:       article.Status,
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}

	// 防御式处理：检查 Author 是否加载成功
	if article.Author.ID > 0 {
		resp.Author = &dto.UserResponse{
			ID:       article.Author.ID,
			Username: article.Author.Username,
			Nickname: article.Author.Nickname,
			Avatar:   article.Author.Avatar,
			Bio:      article.Author.Bio,
			Role:     article.Author.Role,
			Email:    article.Author.Email, // 仅在允许的场景显示
		}
	}

	// 防御式处理：检查 Category 是否加载成功
	if article.Category.ID > 0 {
		resp.Category = &dto.CategoryDTO{
			ID:          article.Category.ID,
			Name:        article.Category.Name,
			Description: article.Category.Description,
			Status:      article.Category.Status,
			CreatedAt:   article.Category.CreatedAt,
			UpdatedAt:   article.Category.UpdatedAt,
		}
	}

	// 防御式处理：Tags 初始化为空切片而非 nil
	if len(article.Tags) > 0 {
		resp.Tags = make([]dto.TagDTO, 0, len(article.Tags))
		for _, tag := range article.Tags {
			resp.Tags = append(resp.Tags, dto.TagDTO{
				ID:        tag.ID,
				Name:      tag.Name,
				Color:     tag.Color,
				Status:    tag.Status,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			})
		}
	} else {
		resp.Tags = []dto.TagDTO{} // 返回空切片而非 nil
	}

	return resp
}
