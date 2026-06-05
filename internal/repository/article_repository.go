package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// withPreloads 统一的预加载方法（加载所有关联）
func (r *ArticleRepository) withPreloads(tx *gorm.DB) *gorm.DB {
	return tx.Preload("Author").
		Preload("Category").
		Preload("Tags")
}

// Create 创建文章
func (r *ArticleRepository) Create(article *entity.Article) error {
	return r.db.Create(article).Error
}

// GetByID 根据 ID 查询文章（带关联）
func (r *ArticleRepository) GetByID(id uint) (*entity.Article, error) {
	var article entity.Article
	err := r.withPreloads(r.db).First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetList 获取文章列表（分页）
func (r *ArticleRepository) GetList(page, pageSize int, status *int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.withPreloads(r.db.Model(&entity.Article{})).
		Order("created_at DESC")

	// 状态筛选（过滤掉草稿和已删除的文章）
	if status != nil {
		query = query.Where("status = ?", *status)
	} else {
		query = query.Where("status = ?", 1)
	}

	// 总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetByAuthorID 根据作者 ID 查询文章
func (r *ArticleRepository) GetByAuthorID(authorID uint, page, pageSize int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.withPreloads(r.db.Model(&entity.Article{})).
		Where("author_id = ?", authorID).
		Order("created_at DESC")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetByCategoryID 根据分类 ID 查询文章
func (r *ArticleRepository) GetByCategoryID(categoryID uint, page, pageSize int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.withPreloads(r.db.Model(&entity.Article{})).
		Where("category_id = ? AND status = ?", categoryID, 1).
		Order("created_at DESC")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetByTagID 根据标签 ID 查询文章
func (r *ArticleRepository) GetByTagID(tagID uint, page, pageSize int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	// 构建基础查询
	baseQuery := r.db.Model(&entity.Article{}).
		Joins("JOIN article_tag ON article_tag.article_id = articles.id").
		Where("article_tag.tag_id = ? AND articles.status = ?", tagID, 1)

	// 获取总数
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询文章ID
	offset := (page - 1) * pageSize
	var articleIDs []uint
	if err := baseQuery.Order("articles.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Pluck("articles.id", &articleIDs).Error; err != nil {
		return nil, 0, err
	}

	if len(articleIDs) == 0 {
		return articles, total, nil
	}

	// 根据ID查询完整的文章信息（包含预加载）
	if err := r.withPreloads(r.db.Where("id IN ?", articleIDs)).
		Order("created_at DESC").
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// Update 更新文章
func (r *ArticleRepository) Update(article *entity.Article) error {
	return r.db.Save(article).Error
}

// Delete 删除文章（软删除）
func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Model(&entity.Article{}).
		Where("id = ?", id).
		Update("status", -1).Error
}

// IncrementViewCount 增加浏览量
func (r *ArticleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&entity.Article{}).
		Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementLikeCount 增加点赞数
func (r *ArticleRepository) IncrementLikeCount(id uint) error {
	return r.db.Model(&entity.Article{}).
		Where("id = ?", id).
		Update("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCount 减少点赞数
func (r *ArticleRepository) DecrementLikeCount(id uint) error {
	return r.db.Model(&entity.Article{}).
		Where("id = ? AND like_count > 0", id).
		Update("like_count", gorm.Expr("like_count - 1")).Error
}

// IncrementCommentCount 增加评论数
func (r *ArticleRepository) IncrementCommentCount(id uint) error {
	return r.db.Model(&entity.Article{}).
		Where("id = ?", id).
		Update("comment_count", gorm.Expr("comment_count + 1")).Error
}

// DecrementCommentCount 减少评论数
func (r *ArticleRepository) DecrementCommentCount(id uint) error {
	return r.db.Model(&entity.Article{}).
		Where("id = ? AND comment_count > 0", id).
		Update("comment_count", gorm.Expr("comment_count - 1")).Error
}
