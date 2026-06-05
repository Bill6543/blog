package api

import (
	"blog/internal/middleware"
	"blog/internal/model/dto"
	"blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentHandler 评论 Handler
type CommentHandler struct {
	*Handler
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	comment, err := h.CommentService.CreateComment(req, userID)
	if err != nil {
		switch {
		case errors.IsArticleNotFound(err):
			response.BadRequest(c, "文章不存在")
		case errors.IsCommentNotFound(err):
			response.BadRequest(c, "父评论不存在")
		default:
			logger.Errorf("Create comment failed: error=%v", err)
			response.InternalError(c, "创建失败")
		}
		return
	}

	response.Success(c, comment)
}

// GetComment 获取评论详情
func (h *CommentHandler) GetComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
		return
	}

	comment, err := h.CommentService.GetCommentByID(uint(id))
	if err != nil {
		switch {
		case errors.IsCommentNotFound(err):
			response.NotFound(c, "评论不存在")
		default:
			logger.Errorf("Get comment failed: error=%v", err)
			response.InternalError(c, "获取评论失败")
		}
		return
	}

	response.Success(c, comment)
}

// GetCommentList 获取评论列表
func (h *CommentHandler) GetCommentList(c *gin.Context) {
	// 获取并验证 article_id 参数
	articleIDStr := c.Query("article_id")
	if articleIDStr == "" {
		response.BadRequest(c, "缺少 article_id 参数")
		return
	}

	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 article_id")
		return
	}

	list, err := h.CommentService.GetCommentsByArticleID(uint(articleID))
	if err != nil {
		logger.Errorf("Get comment list failed: error=%v", err)
		response.InternalError(c, "获取评论列表失败")
		return
	}

	response.Success(c, list)
}

// UpdateComment 更新评论
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
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

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	comment, err := h.CommentService.UpdateComment(uint(id), req, userID, isAdmin)
	if err != nil {
		switch {
		case errors.IsCommentNotFound(err):
			response.NotFound(c, "评论不存在")
		case errors.IsPermissionDenied(err):
			response.Forbidden(c, "无权修改此评论")
		default:
			logger.Errorf("Update comment failed: error=%v", err)
			response.InternalError(c, "更新失败")
		}
		return
	}

	response.Success(c, comment)
}

// DeleteComment 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
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

	if err := h.CommentService.DeleteComment(uint(id), userID, isAdmin); err != nil {
		switch {
		case errors.IsCommentNotFound(err):
			response.NotFound(c, "评论不存在")
		case errors.IsPermissionDenied(err):
			response.Forbidden(c, "无权删除此评论")
		default:
			logger.Errorf("Delete comment failed: error=%v", err)
			response.InternalError(c, "删除失败")
		}
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetPendingComments 获取待审核评论列表（管理员专用）
func (h *CommentHandler) GetPendingComments(c *gin.Context) {
	comments, err := h.CommentService.GetPendingComments()
	if err != nil {
		logger.Errorf("Get pending comments failed: error=%v", err)
		response.InternalError(c, "获取待审核评论失败")
		return
	}

	response.Success(c, comments)
}

// GetPendingCommentCount 获取待审核评论数量（管理员专用）
func (h *CommentHandler) GetPendingCommentCount(c *gin.Context) {
	count, err := h.CommentService.GetPendingCommentCount()
	if err != nil {
		logger.Errorf("Get pending comment count failed: error=%v", err)
		response.InternalError(c, "获取待审核评论数量失败")
		return
	}

	response.Success(c, count)
}

// ApproveComment 审核通过评论（管理员专用）
func (h *CommentHandler) ApproveComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
		return
	}

	if err := h.CommentService.ApproveComment(uint(id)); err != nil {
		switch {
		case errors.IsCommentNotFound(err):
			response.NotFound(c, "评论不存在")
		default:
			logger.Errorf("Approve comment failed: error=%v", err)
			response.InternalError(c, "审核通过失败")
		}
		return
	}

	response.SuccessWithMessage(c, "评论审核通过", nil)
}

// RejectComment 审核拒绝评论（管理员专用）
func (h *CommentHandler) RejectComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
		return
	}

	if err := h.CommentService.RejectComment(uint(id)); err != nil {
		switch {
		case errors.IsCommentNotFound(err):
			response.NotFound(c, "评论不存在")
		default:
			logger.Errorf("Reject comment failed: error=%v", err)
			response.InternalError(c, "审核拒绝失败")
		}
		return
	}

	response.SuccessWithMessage(c, "评论审核拒绝", nil)
}
