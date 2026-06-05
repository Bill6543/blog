package service

import (
	"blog/internal/model/dto"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/logger"
)

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}

// CreateTag 创建标签
func (s *TagService) CreateTag(req dto.CreateTagRequest) (*dto.TagDTO, error) {
	// 检查名称是否存在
	if s.tagRepo.ExistsByName(req.Name) {
		logger.Warnf("Create tag failed: name already exists, name=%s", req.Name)
		return nil, errors.New(errors.TagNameExistsCode)
	}

	tag := &entity.Tag{
		Name:   req.Name,
		Color:  req.Color,
		Status: 1,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		logger.Errorf("Create tag failed: database error, name=%s, error=%v", req.Name, err)
		return nil, err
	}

	return s.convertToDTO(tag), nil
}

// GetAllTags 获取所有标签
func (s *TagService) GetAllTags() ([]dto.TagDTO, error) {
	tags, err := s.tagRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return s.convertToDTOs(tags), nil
}

// GetTagByID 根据 ID 查询标签
func (s *TagService) GetTagByID(id uint) (*dto.TagDTO, error) {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		logger.Warnf("Get tag failed: tag not found, tagID=%d", id)
		return nil, errors.New(errors.TagNotFoundCode)
	}

	return s.convertToDTO(tag), nil
}

// UpdateTag 更新标签
func (s *TagService) UpdateTag(id uint, req dto.UpdateTagRequest) (*dto.TagDTO, error) {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		logger.Warnf("Update tag failed: tag not found, tagID=%d", id)
		return nil, errors.New(errors.TagNotFoundCode)
	}

	// 更新字段
	if req.Name != "" && req.Name != tag.Name {
		// Only check name uniqueness if name is actually changed
		if s.tagRepo.ExistsByName(req.Name) {
			logger.Warnf("Update tag failed: name already exists, name=%s", req.Name)
			return nil, errors.New(errors.TagNameExistsCode)
		}
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}

	if err := s.tagRepo.Update(tag); err != nil {
		logger.Errorf("Update tag failed: database error, tagID=%d, error=%v", id, err)
		return nil, err
	}

	return s.convertToDTO(tag), nil
}

// DeleteTag 删除标签
func (s *TagService) DeleteTag(id uint) error {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		logger.Warnf("Delete tag failed: tag not found, tagID=%d", id)
		return errors.New(errors.TagNotFoundCode)
	}

	if err := s.tagRepo.Delete(id); err != nil {
		logger.Errorf("Delete tag failed: database error, tagID=%d, name=%s, error=%v", id, tag.Name, err)
		return err
	}

	return nil
}

// convertToDTO 转换为 DTO
func (s *TagService) convertToDTO(tag *entity.Tag) *dto.TagDTO {
	return &dto.TagDTO{
		ID:        tag.ID,
		Name:      tag.Name,
		Color:     tag.Color,
		Status:    tag.Status,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}

// convertToDTOs 批量转换
func (s *TagService) convertToDTOs(tags []entity.Tag) []dto.TagDTO {
	var dtos []dto.TagDTO
	for _, tag := range tags {
		dto := s.convertToDTO(&tag)
		dtos = append(dtos, *dto)
	}
	return dtos
}
