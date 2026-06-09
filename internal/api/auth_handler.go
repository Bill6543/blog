package api

import (
	"blog/internal/middleware"
	"blog/internal/model/dto"
	"blog/internal/service"
	"blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证 Handler
type AuthHandler struct {
	*Handler
}

// Register 用户注册 (支持表单和头像上传)
func (h *AuthHandler) Register(c *gin.Context) {
	// 1. 自动绑定并验证（支持 multipart/form-data）
	var req dto.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		// 记录详细错误到日志
		logger.Errorf("Register parameter validation failed: %v", err)
		// 返回友好的错误消息
		response.BadRequest(c, errors.HandleValidationError(err))
		return
	}

	// 2. 处理头像上传
	var avatarPath string = h.UploadService.GetDefaultAvatar()
	file, err := c.FormFile("avatar")
	if err == nil && file != nil {
		logger.Infof("Avatar upload attempt: filename=%s, size=%d, content_type=%s", file.Filename, file.Size, file.Header.Get("Content-Type"))
		avatarPath, err = h.UploadService.UploadImage(file, service.UploadConfig{
			Dir:     service.ImageTypeAvatar,
			MaxSize: 2 * 1024 * 1024,
		})
		if err != nil {
			logger.Errorf("Avatar upload failed: filename=%s, error=%v", file.Filename, err)
			response.BadRequest(c, "头像上传失败: "+err.Error())
			return
		}
		logger.Infof("Avatar uploaded successfully: path=%s", avatarPath)
	}

	// 3. 调用 Service
	user, err := h.AuthService.Register(req, avatarPath)
	if err != nil {
		// 2. Service 层返回英文错误，转换为中文返回给用户
		switch {
		case errors.IsUsernameExists(err):
			response.BadRequest(c, "用户名已存在")
		case errors.IsEmailExists(err):
			response.BadRequest(c, "邮箱已被注册")
		case errors.IsPasswordHashFailed(err):
			response.InternalError(c, "密码加密失败")
		case errors.IsRegistrationFailed(err):
			response.InternalError(c, "注册失败，请稍后重试")
		default:
			// 其他错误（如密码强度不足等）直接返回
			logger.Errorf("Register service error: %v", err)
			response.BadRequest(c, err.Error())
		}
		return
	}

	response.Success(c, user)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 记录详细错误到日志
		logger.Errorf("Login parameter validation failed: %v", err)
		// 返回友好的错误消息
		response.BadRequest(c, errors.HandleValidationError(err))
		return
	}

	token, err := h.AuthService.Login(req)
	if err != nil {
		// 2. 根据错误类型返回不同的响应
		switch {
		case errors.IsInvalidCredentials(err):
			response.Unauthorized(c, "用户名或密码错误")
		case errors.IsUserDisabled(err):
			response.Forbidden(c, "用户已被禁用")
		case errors.IsUpdateSessionFailed(err), errors.IsGenerateTokenFailed(err):
			response.InternalError(c, "登录失败，请稍后重试")
		default:
			response.Unauthorized(c, "用户名或密码错误")
		}
		return
	}

	response.Success(c, gin.H{"token": token})
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取自己的完整信息（包含 Email）
	user, err := h.UserService.GetUserInfoForSelf(userID)
	if err != nil {
		response.InternalError(c, "获取用户信息失败")
		return
	}

	response.Success(c, user)
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	if err := h.AuthService.Logout(userID); err != nil {
		logger.Errorf("Logout failed: error=%v", err)
		response.InternalError(c, "登出失败，请稍后重试")
		return
	}

	response.SuccessWithMessage(c, "登出成功", nil)
}

// ForgotPassword 忘记密码
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("ForgotPassword parameter validation failed: %v", err)
		response.BadRequest(c, errors.HandleValidationError(err))
		return
	}

	resetToken, err := h.AuthService.ForgotPassword(req.Email)
	if err != nil {
		// 统一错误提示，防止邮箱枚举攻击
		switch {
		case errors.IsUserNotFound(err):
			logger.Warnf("Forgot password: user not found, email=%s", req.Email)
		default:
			logger.Errorf("ForgotPassword service error: %v", err)
		}
		response.SuccessWithMessage(c, "如果该邮箱已注册，重置链接将发送到您的邮箱", nil)
		return
	}

	// 构建重置链接
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)

	response.SuccessWithMessage(c, "重置链接已生成", map[string]interface{}{
		"reset_url":  resetURL,
		"expires_in": 3600, // 1小时
	})
}

// ResetPassword 重置密码
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("ResetPassword parameter validation failed: %v", err)
		response.BadRequest(c, errors.HandleValidationError(err))
		return
	}

	if err := h.AuthService.ResetPassword(req.Token, req.NewPassword); err != nil {
		switch {
		case errors.IsInvalidToken(err):
			response.BadRequest(c, "无效的重置令牌")
		default:
			logger.Errorf("ResetPassword service error: %v", err)
			response.InternalError(c, err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "密码重置成功，请使用新密码登录", nil)
}
