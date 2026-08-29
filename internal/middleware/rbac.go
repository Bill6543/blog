package middleware

import (
	"blog/pkg/logger"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// RoleGuard 角色权限校验中间件
// 用于校验用户是否具有指定的角色
func RoleGuard(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			logger.Warnf("Role guard failed: role not found in context, path=%s", c.Request.URL.Path)
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			logger.Errorf("Role guard failed: invalid role type %T in context, path=%s", role, c.Request.URL.Path)
			response.Unauthorized(c, "无效的用户角色")
			c.Abort()
			return
		}

		// 检查是否在允许的角色列表中
		for _, required := range requiredRoles {
			if roleStr == required {
				c.Next()
				return
			}
		}

		// 权限拒绝，记录日志
		logger.Warnf("Permission denied: user=%s, role=%s, path=%s, required_roles=%v",
			GetUsernameFromContext(c), roleStr, c.Request.URL.Path, requiredRoles)
		response.Forbidden(c, "权限不足")
		c.Abort()
	}
}

// IsAdmin 快捷方法：仅允许管理员
func IsAdmin() gin.HandlerFunc {
	return RoleGuard("admin")
}

// 注：资源所有权（如"只能改自己的文章/评论"）的校验放在 service 层，
// 因为需要先查库拿到资源的 owner 才能判断，中间件层拿不到。
// 参见 article_service.go / comment_service.go 中的 PermissionDenied 分支。
