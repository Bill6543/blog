package service

import (
	"blog/internal/model/dto"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/logger"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(req dto.CreateCategoryRequest) (*dto.CategoryDTO, error) {
	// 检查名称是否存在
	if s.categoryRepo.ExistsByName(req.Name) {
		logger.Warnf("Create category failed: name already exists, name=%s", req.Name)
		return nil, errors.New(errors.CategoryExistsCode)
	}

	category := &entity.Category{
		Name:        req.Name,
		Description: req.Description,
		Status:      1,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		logger.Errorf("Create category failed: database error, name=%s, error=%v", req.Name, err)
		return nil, err
	}

	return s.convertToDTO(category), nil
}

// GetAllCategories 获取所有分类
func (s *CategoryService) GetAllCategories() ([]dto.CategoryDTO, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return s.convertToDTOs(categories), nil
}

// GetCategoryByID 获取分类详情
func (s *CategoryService) GetCategoryByID(id uint) (*dto.CategoryDTO, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		logger.Warnf("Get category failed: category not found, categoryID=%d", id)
		return nil, errors.New(errors.CategoryNotFoundCode)
	}

	return s.convertToDTO(category), nil
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(id uint, req dto.UpdateCategoryRequest) (*dto.CategoryDTO, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New(errors.CategoryNotFoundCode)
	}

	// 更新字段
	if req.Name != "" {
		if s.categoryRepo.ExistsByName(req.Name) {
			return nil, errors.New(errors.CategoryExistsCode)
		}
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return s.convertToDTO(category), nil
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(id uint) error {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		logger.Warnf("Delete category failed: category not found, categoryID=%d", id)
		return errors.New(errors.CategoryNotFoundCode)
	}

	// 检查分类是否被文章使用
	if s.categoryRepo.IsUsedByArticles(id) {
		logger.Warnf("Delete category failed: category is in use by articles, categoryID=%d, name=%s", id, category.Name)
		return errors.New(errors.CategoryInUseCode)
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		logger.Errorf("Delete category failed: database error, categoryID=%d, name=%s, error=%v", id, category.Name, err)
		return err
	}

	return nil
}

// convertToDTO 转换为 DTO
func (s *CategoryService) convertToDTO(cat *entity.Category) *dto.CategoryDTO {
	return &dto.CategoryDTO{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		Status:      cat.Status,
		CreatedAt:   cat.CreatedAt,
		UpdatedAt:   cat.UpdatedAt,
	}
}

// convertToDTOs 批量转换
func (s *CategoryService) convertToDTOs(categories []entity.Category) []dto.CategoryDTO {
	var dtos []dto.CategoryDTO
	for _, cat := range categories {
		dtos = append(dtos, *s.convertToDTO(&cat))
	}
	return dtos
}
