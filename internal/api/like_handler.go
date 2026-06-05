package api

import (
	"blog/internal/middleware"
	"blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// LikeHandler 点赞 Handler
type LikeHandler struct {
	*Handler
}

// LikeArticle 点赞文章
func (h *LikeHandler) LikeArticle(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	if err := h.LikeService.LikeArticle(uint(articleID), userID); err != nil {
		switch {
		case errors.IsArticleNotFound(err):
			response.NotFound(c, "文章不存在")
		case errors.IsUserNotFound(err):
			response.BadRequest(c, "用户不存在")
		case errors.IsAlreadyLiked(err):
			response.BadRequest(c, "已点赞")
		default:
			logger.Errorf("Like article failed: error=%v", err)
			response.InternalError(c, "点赞失败")
		}
		return
	}

	response.SuccessWithMessage(c, "点赞成功", nil)
}

// UnlikeArticle 取消点赞
func (h *LikeHandler) UnlikeArticle(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	if err := h.LikeService.UnlikeArticle(uint(articleID), userID); err != nil {
		switch {
		case errors.IsNotLikedYet(err):
			response.BadRequest(c, "尚未点赞")
		default:
			logger.Errorf("Unlike article failed: error=%v", err)
			response.InternalError(c, "取消点赞失败")
		}
		return
	}

	response.SuccessWithMessage(c, "取消点赞成功", nil)
}

// GetLikeStatus 获取点赞状态
func (h *LikeHandler) GetLikeStatus(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	liked, err := h.LikeService.GetLikeStatus(uint(articleID), userID)
	if err != nil {
		logger.Errorf("Get like status failed: error=%v", err)
		response.InternalError(c, "获取点赞状态失败")
		return
	}

	response.Success(c, gin.H{
		"liked": liked,
	})
}
