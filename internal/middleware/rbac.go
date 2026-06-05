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

// IsAdminOrOwner 管理员或资源所有者校验中间件
// 需要配合上下文中的 user_id 使用
func IsAdminOrOwner(resourceOwnerID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserIDFromContext(c)
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		role, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			logger.Errorf("Invalid role type in context: %T", role)
			response.Unauthorized(c, "无效的用户角色")
			c.Abort()
			return
		}

		// 管理员直接放行
		if roleStr == "admin" {
			c.Next()
			return
		}

		// 检查是否为资源所有者
		if userID != resourceOwnerID {
			logger.Warnf("Permission denied: user=%d is not owner (owner=%d), path=%s",
				userID, resourceOwnerID, c.Request.URL.Path)
			response.Forbidden(c, "无权操作他人资源")
			c.Abort()
			return
		}

		c.Next()
	}
}
