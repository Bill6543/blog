package api

import (
	"blog/internal/model/dto"
	"blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// TagHandler 标签 Handler
type TagHandler struct {
	*Handler
}

// CreateTag 创建标签
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag, err := h.TagService.CreateTag(req)
	if err != nil {
		switch {
		case errors.IsTagNameExists(err):
			response.BadRequest(c, "标签名称已存在")
		default:
			logger.Errorf("Create tag failed: error=%v", err)
			response.InternalError(c, "创建失败")
		}
		return
	}

	response.Success(c, tag)
}

// GetTag 获取标签详情
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的标签 ID")
		return
	}

	tag, err := h.TagService.GetTagByID(uint(id))
	if err != nil {
		switch {
		case errors.IsTagNotFound(err):
			response.NotFound(c, "标签不存在")
		default:
			logger.Errorf("Get tag failed: error=%v", err)
			response.InternalError(c, "获取失败")
		}
		return
	}

	response.Success(c, tag)
}

// GetTagList 获取标签列表
func (h *TagHandler) GetTagList(c *gin.Context) {
	list, err := h.TagService.GetAllTags()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, list)
}

// UpdateTag 更新标签
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的标签 ID")
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag, err := h.TagService.UpdateTag(uint(id), req)
	if err != nil {
		switch {
		case errors.IsTagNotFound(err):
			response.NotFound(c, "标签不存在")
		case errors.IsTagNameExists(err):
			response.BadRequest(c, "标签名称已存在")
		default:
			logger.Errorf("Update tag failed: error=%v", err)
			response.InternalError(c, "更新失败")
		}
		return
	}

	response.Success(c, tag)
}

// DeleteTag 删除标签
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的标签 ID")
		return
	}

	if err := h.TagService.DeleteTag(uint(id)); err != nil {
		switch {
		case errors.IsTagNotFound(err):
			response.NotFound(c, "标签不存在")
		default:
			logger.Errorf("Delete tag failed: error=%v", err)
			response.InternalError(c, "删除失败")
		}
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}


