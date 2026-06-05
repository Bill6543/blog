package api

import (
	"blog/internal/model/dto"
	"blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类 Handler
type CategoryHandler struct {
	*Handler
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := h.CategoryService.CreateCategory(req)
	if err != nil {
		switch {
		case errors.IsCategoryExists(err):
			response.BadRequest(c, "分类名称已存在")
		default:
			logger.Errorf("Create category failed: error=%v", err)
			response.InternalError(c, "创建失败")
		}
		return
	}

	response.Success(c, category)
}

// GetCategory 获取分类详情
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类 ID")
		return
	}

	dtoObj, err := h.CategoryService.GetCategoryByID(uint(id))
	if err != nil {
		switch {
		case errors.IsCategoryNotFound(err):
			response.NotFound(c, "分类不存在")
		default:
			logger.Errorf("Get category failed: error=%v", err)
			response.InternalError(c, "获取分类失败")
		}
		return
	}

	response.Success(c, dtoObj)
}

// GetCategoryList 获取分类列表
func (h *CategoryHandler) GetCategoryList(c *gin.Context) {
	list, err := h.CategoryService.GetAllCategories()
	if err != nil {
		logger.Errorf("Get category list failed: error=%v", err)
		response.InternalError(c, "获取分类列表失败")
		return
	}

	response.Success(c, list)
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类 ID")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := h.CategoryService.UpdateCategory(uint(id), req)
	if err != nil {
		switch {
		case errors.IsCategoryNotFound(err):
			response.NotFound(c, "分类不存在")
		case errors.IsCategoryExists(err):
			response.BadRequest(c, "分类名称已存在")
		default:
			logger.Errorf("Update category failed: error=%v", err)
			response.InternalError(c, "更新失败")
		}
		return
	}

	response.Success(c, category)
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类 ID")
		return
	}

	if err := h.CategoryService.DeleteCategory(uint(id)); err != nil {
		switch {
		case errors.IsCategoryNotFound(err):
			response.NotFound(c, "分类不存在")
		case errors.IsCategoryInUse(err):
			response.BadRequest(c, "该分类下有文章，无法删除")
		default:
			logger.Errorf("Delete category failed: error=%v", err)
			response.InternalError(c, "删除失败")
		}
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
