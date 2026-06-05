package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create 创建标签
func (r *TagRepository) Create(tag *entity.Tag) error {
	return r.db.Create(tag).Error
}

// GetByID 根据 ID 查询标签
func (r *TagRepository) GetByID(id uint) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetAll 获取所有标签
func (r *TagRepository) GetAll() ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.Where("status = ?", 1).
		Order("created_at DESC").
		Find(&tags).Error
	return tags, err
}

// Update 更新标签
func (r *TagRepository) Update(tag *entity.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签
func (r *TagRepository) Delete(id uint) error {
	return r.db.Model(&entity.Tag{}).
		Where("id = ?", id).
		Update("status", 0).Error
}

// ExistsByName 检查标签名称是否存在
func (r *TagRepository) ExistsByName(name string) bool {
	var count int64
	r.db.Model(&entity.Tag{}).Where("name = ?", name).Count(&count)
	return count > 0
}

// GetByIDs 根据 ID 数组查询标签
func (r *TagRepository) GetByIDs(ids []uint) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}
