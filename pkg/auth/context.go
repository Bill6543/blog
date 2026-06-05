package auth

import (
	"blog/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Role 用户角色
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleUser     Role = "user"
	RoleAnonymous Role = ""
)

// GetRoleFromContext 安全地从 Context 获取角色
// 如果角色不存在或类型错误，返回 RoleAnonymous
func GetRoleFromContext(c *gin.Context) Role {
	role, exists := c.Get("role")
	if !exists {
		logger.Debugf("Role not set in context, returning RoleAnonymous")
		return RoleAnonymous
	}

	roleStr, ok := role.(string)
	if !ok {
		logger.Errorf("Invalid role type in context: %T", role)
		return RoleAnonymous
	}

	// 转换为 Role 类型
	return Role(roleStr)
}

// IsAdmin 检查是否是管理员
// 如果角色获取失败，返回 false
func IsAdmin(c *gin.Context) bool {
	return GetRoleFromContext(c) == RoleAdmin
}

// IsUser 检查是否是普通用户
// 如果角色获取失败，返回 false
func IsUser(c *gin.Context) bool {
	return GetRoleFromContext(c) == RoleUser
}

// GetRoleOrDefault 获取角色，如果失败则返回默认值
func GetRoleOrDefault(c *gin.Context, defaultRole Role) Role {
	role := GetRoleFromContext(c)
	if role == RoleAnonymous {
		return defaultRole
	}
	return role
}

// IsAuthenticated 检查是否已登录（通过 user_id 判断）
func IsAuthenticated(c *gin.Context) bool {
	userID, exists := c.Get("user_id")
	if !exists {
		return false
	}

	id, ok := userID.(uint)
	return ok && id != 0
}
