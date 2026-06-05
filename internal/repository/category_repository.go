package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create 创建分类
func (r *CategoryRepository) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

// GetByID 根据 ID 查询分类
func (r *CategoryRepository) GetByID(id uint) (*entity.Category, error) {
	var category entity.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll 获取所有分类
func (r *CategoryRepository) GetAll() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Where("status = ?", 1).
		Order("created_at DESC").
		Find(&categories).Error
	return categories, err
}

// Update 更新分类
func (r *CategoryRepository) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类（软删除）
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Model(&entity.Category{}).
		Where("id = ?", id).
		Update("status", 0).Error
}

// ExistsByName 检查分类名称是否存在
func (r *CategoryRepository) ExistsByName(name string) bool {
	var count int64
	r.db.Model(&entity.Category{}).Where("name = ?", name).Count(&count)
	return count > 0
}

// IsUsedByArticles 检查分类是否被文章使用
func (r *CategoryRepository) IsUsedByArticles(categoryID uint) bool {
	var count int64
	r.db.Model(&entity.Article{}).Where("category_id = ? AND status != -1", categoryID).Count(&count)
	return count > 0
}
