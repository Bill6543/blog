package api

import (
	"blog/internal/middleware"
	"blog/internal/model/dto"
	"blog/internal/service"
	"blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ArticleHandler 文章 Handler
type ArticleHandler struct {
	*Handler
}

// CreateArticle 创建文章
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.CreateArticleRequest

	// 1. 自动绑定并验证（支持 multipart/form-data）
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "参数验证失败："+err.Error())
		return
	}

	// 2. 处理封面上传（可选）
	file, err := c.FormFile("cover_image_file")
	if err == nil && file != nil {
		coverPath, err := h.UploadService.UploadImage(file, service.UploadConfig{
			Dir:     service.ImageTypeCover,
			MaxSize: 5 * 1024 * 1024,
		})
		if err != nil {
			response.BadRequest(c, "封面上传失败")
			return
		}
		req.CoverImage = coverPath
	}

	// 3. 调用 Service 创建文章（已返回 DTO）
	article, err := h.ArticleService.CreateArticle(req, userID)
	if err != nil {
		switch {
		case errors.IsCategoryNotFound(err):
			response.BadRequest(c, "分类不存在")
		case errors.IsTagNotFound(err):
			response.BadRequest(c, "标签不存在")
		default:
			logger.Errorf("Create article failed: error=%v", err)
			response.InternalError(c, "创建失败")
		}
		return
	}

	response.Success(c, article)
}

// GetArticle 获取文章详情
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	article, err := h.ArticleService.GetArticleByID(uint(id))
	if err != nil {
		switch {
		case errors.IsArticleNotFound(err):
			response.NotFound(c, "文章不存在")
		default:
			logger.Errorf("Get article failed: error=%v", err)
			response.InternalError(c, "获取文章失败")
		}
		return
	}

	response.Success(c, article)
}

// GetArticleList 获取文章列表
func (h *ArticleHandler) GetArticleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int
	if s := c.Query("status"); s != "" {
		statusVal, _ := strconv.Atoi(s)
		status = &statusVal
	}

	var categoryID *uint
	if catID := c.Query("category_id"); catID != "" {
		catIDVal, _ := strconv.ParseUint(catID, 10, 32)
		catIDUint := uint(catIDVal)
		categoryID = &catIDUint
	}

	var tagID *uint
	if tID := c.Query("tag_id"); tID != "" {
		tIDVal, _ := strconv.ParseUint(tID, 10, 32)
		tIDUint := uint(tIDVal)
		tagID = &tIDUint
	}

	list, err := h.ArticleService.GetArticleList(page, pageSize, status, categoryID, tagID)
	if err != nil {
		logger.Errorf("Get article list failed: error=%v", err)
		response.InternalError(c, "获取文章列表失败")
		return
	}

	response.Success(c, list)
}

// UpdateArticle 更新文章
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 判断是否为管理员
	role, _ := c.Get("role")
	isAdmin := role == "admin"

	var req dto.UpdateArticleRequest

	// 1. 自动绑定并验证（支持 multipart/form-data）
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "参数验证失败："+err.Error())
		return
	}

	// 2. 处理封面上传（可选）
	file, err := c.FormFile("cover_image_file")
	if err == nil && file != nil {
		coverPath, err := h.UploadService.UploadImage(file, service.UploadConfig{
			Dir:     service.ImageTypeCover,
			MaxSize: 5 * 1024 * 1024,
		})
		if err != nil {
			response.BadRequest(c, "封面上传失败")
			return
		}
		req.CoverImage = coverPath
	}

	// 3. 调用 Service 更新文章
	article, err := h.ArticleService.UpdateArticle(uint(id), req, userID, isAdmin)
	if err != nil {
		switch {
		case errors.IsArticleNotFound(err):
			response.NotFound(c, "文章不存在")
		case errors.IsPermissionDenied(err):
			response.Forbidden(c, "无权修改此文章")
		default:
			logger.Errorf("Update article failed: error=%v", err)
			response.InternalError(c, "更新失败")
		}
		return
	}

	response.Success(c, article)
}

// DeleteArticle 删除文章
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 判断是否为管理员
	role, _ := c.Get("role")
	isAdmin := role == "admin"

	if err := h.ArticleService.DeleteArticle(uint(id), userID, isAdmin); err != nil {
		switch {
		case errors.IsArticleNotFound(err):
			response.NotFound(c, "文章不存在")
		case errors.IsPermissionDenied(err):
			response.Forbidden(c, "无权删除此文章")
		default:
			logger.Errorf("Delete article failed: error=%v", err)
			response.InternalError(c, "删除失败")
		}
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetUserArticles 获取指定用户的文章列表
func (h *ArticleHandler) GetUserArticles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户 ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int
	if s := c.Query("status"); s != "" {
		statusVal, _ := strconv.Atoi(s)
		status = &statusVal
	}

	list, err := h.ArticleService.GetUserArticles(uint(userID), page, pageSize, status)
	if err != nil {
		logger.Errorf("Get user articles failed: error=%v", err)
		response.InternalError(c, "获取用户文章列表失败")
		return
	}

	response.Success(c, list)
}

// RestoreArticle 恢复已删除的文章
func (h *ArticleHandler) RestoreArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 判断是否为管理员
	role, _ := c.Get("role")
	isAdmin := role == "admin"

	if err := h.ArticleService.RestoreArticle(uint(id), userID, isAdmin); err != nil {
		switch {
		case errors.IsArticleNotFound(err):
			response.NotFound(c, "文章不存在")
		case errors.IsPermissionDenied(err):
			response.Forbidden(c, "无权恢复此文章")
		default:
			logger.Errorf("Restore article failed: error=%v", err)
			response.InternalError(c, "恢复失败")
		}
		return
	}

	response.SuccessWithMessage(c, "恢复成功", nil)
}

// IncreaseViewCount 增加文章浏览量
func (h *ArticleHandler) IncreaseViewCount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	if err := h.ArticleService.IncreaseViewCount(uint(id)); err != nil {
		logger.Errorf("Increase view count failed: error=%v", err)
		response.InternalError(c, "增加浏览量失败")
		return
	}

	response.SuccessWithMessage(c, "浏览量增加成功", nil)
}
