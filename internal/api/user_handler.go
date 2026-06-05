package api

import (
	"blog/internal/middleware"
	"blog/internal/model/dto"
	"blog/pkg/auth"
	"blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// UserHandler 用户 Handler
type UserHandler struct {
	*Handler
}

// GetUserInfo 获取用户信息（管理员看完整，普通用户看公开）
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户 ID")
		return
	}

	// ✅ 安全获取角色
	viewerRole := auth.GetRoleFromContext(c)

	viewerUserID := middleware.GetUserIDFromContext(c)

	var user interface{}

	// 根据角色调用不同的 Service 方法
	if viewerRole == auth.RoleAdmin {
		// 管理员查看完整信息
		user, err = h.UserService.GetUserInfoForAdmin(uint(id))
	} else if viewerUserID != 0 && viewerUserID == uint(id) {
		// ✅ 查看自己的信息（确保已登录）
		user, err = h.UserService.GetUserInfoForSelf(viewerUserID)
	} else {
		// 查看他人信息（不包含 Email）
		user, err = h.UserService.GetUserInfoForOther(uint(id))
	}

	if err != nil {
		switch {
		case errors.IsUserNotFound(err):
			response.NotFound(c, "用户不存在")
		default:
			response.InternalError(c, "获取用户信息失败")
		}
		return
	}

	response.Success(c, user)
}

// GetUserList 获取用户列表（公开接口）
func (h *UserHandler) GetUserList(c *gin.Context) {
	var req dto.UserListRequest
	var err error
	if err = c.ShouldBindQuery(&req); err != nil {
		// 记录详细错误到日志
		logger.Errorf("User list parameter validation failed: %v", err)
		// 返回友好的错误消息
		response.BadRequest(c, errors.HandleValidationError(err))
		return
	}

	// ✅ 安全获取角色
	viewerRole := auth.GetRoleFromContext(c)

	var list interface{}

	// 根据角色调用不同的 Service 方法
	if viewerRole == auth.RoleAdmin {
		// 管理员查看完整列表（包含 Email）
		list, err = h.UserService.GetUserListForAdmin(req)
	} else {
		// 普通用户/匿名用户查看脱敏列表
		list, err = h.UserService.GetPublicUserList(req)
	}

	if err != nil {
		response.InternalError(c, "获取用户列表失败")
		return
	}

	response.Success(c, list)
}

// UpdateUserInfo 更新用户信息（智能判断：管理员可改角色/状态，普通用户只能改资料）
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户 ID")
		return
	}

	viewerUserID := middleware.GetUserIDFromContext(c)
	if viewerUserID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 安全获取角色
	viewerRole := auth.GetRoleFromContext(c)

	var req dto.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 记录详细错误到日志
		logger.Errorf("Update user parameter validation failed: %v", err)
		// 返回友好的错误消息
		response.BadRequest(c, errors.HandleValidationError(err))
		return
	}

	user, err := h.UserService.UpdateUserInfo(uint(id), req, viewerUserID, string(viewerRole))
	if err != nil {
		switch {
		case errors.IsPermissionDenied(err):
			response.Forbidden(c, "无权修改他人信息")
		case errors.IsUserNotFound(err):
			response.NotFound(c, "用户不存在")
		default:
			response.InternalError(c, "更新失败")
		}
		return
	}

	response.Success(c, user)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 记录详细错误到日志
		logger.Errorf("Change password parameter validation failed: %v", err)
		// 返回友好的错误消息
		response.BadRequest(c, errors.HandleValidationError(err))
		return
	}

	if err := h.UserService.ChangePassword(userID, req); err != nil {
		switch {
		case errors.IsUserNotFound(err):
			response.NotFound(c, "用户不存在")
		case errors.IsOldPasswordWrong(err):
			response.BadRequest(c, "原密码错误")
		default:
			response.InternalError(c, "修改密码失败")
		}
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// DeleteUser 删除用户（仅管理员）
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户 ID")
		return
	}

	// ✅ 安全获取角色
	role := auth.GetRoleFromContext(c)
	if role != auth.RoleAdmin {
		response.Forbidden(c, "权限不足")
		return
	}

	if err := h.UserService.DeleteUser(uint(id), string(role)); err != nil {
		switch {
		case errors.IsUserNotFound(err):
			response.NotFound(c, "用户不存在")
		default:
			response.InternalError(c, "删除失败")
		}
		return
	}

	response.SuccessWithMessage(c, "用户已删除", nil)
}
