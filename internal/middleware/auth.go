package middleware

import (
	"blog/internal/cache"
	"blog/internal/repository"
	"blog/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"strings"

	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthConfig JWT 认证配置
type AuthConfig struct {
	Secret       string
	UserRepo     *repository.UserRepository //用于查询最新 session（缓存未命中时回源）
	SessionCache *cache.SessionCache        //用于缓存 session，减少数据库查询
}

// Auth JWT 认证中间件
func Auth(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 提取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warnf("JWT auth failed: missing authorization header, path=%s", c.Request.URL.Path)
			response.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// 解析 Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			logger.Warnf("JWT auth failed: invalid token format, path=%s", c.Request.URL.Path)
			response.Unauthorized(c, "Invalid token format")
			c.Abort()
			return
		}

		// 验证 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Secret), nil
		})

		if err != nil || !token.Valid {
			logger.Warnf("JWT auth failed: invalid token, path=%s, error=%v", c.Request.URL.Path, err)
			response.Unauthorized(c, "Invalid credentials")
			c.Abort()
			return
		}

		// 提取 Claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 安全提取 user_id
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				logger.Warnf("JWT auth failed: missing user_id claim, path=%s", c.Request.URL.Path)
				response.Unauthorized(c, "Invalid credentials")
				c.Abort()
				return
			}
			userID := uint(userIDFloat)

			// 安全提取 login_session
			tokenSession, ok := claims["login_session"].(string)
			if !ok {
				logger.Warnf("JWT auth failed: missing login_session claim, path=%s", c.Request.URL.Path)
				response.Unauthorized(c, "Invalid credentials")
				c.Abort()
				return
			}

			// 提取用户名/角色（用于上下文与日志）
			username, _ := claims["username"].(string)
			role, _ := claims["role"].(string)

			// 获取最新 session：优先读缓存，未命中时回源数据库并回填
			latestSession := ""
			sessionCached := false
			if config.SessionCache != nil {
				if s, hit, err := config.SessionCache.Get(userID); err == nil && hit {
					latestSession = s
					sessionCached = true
				}
			}

			if !sessionCached {
				user, err := config.UserRepo.GetByID(userID)
				if err != nil {
					logger.Warnf("JWT auth failed: user not found, userID=%d, path=%s", userID, c.Request.URL.Path)
					response.Unauthorized(c, "Invalid credentials")
					c.Abort()
					return
				}
				latestSession = user.LoginSession
				if config.SessionCache != nil {
					_ = config.SessionCache.Set(userID, latestSession)
				}
			}

			// 对比 session
			if latestSession != tokenSession {
				logger.Warnf("JWT auth failed: session mismatch, userID=%d, username=%s, path=%s",
					userID, username, c.Request.URL.Path)
				response.Unauthorized(c, "Your account has logged in on another device, please login again")
				c.Abort()
				return
			}

			// 将用户信息存入上下文
			c.Set("user_id", userID)
			c.Set("username", username)
			c.Set("role", role)
			c.Next()
		}
	}
}

// GetUserIDFromContext 从上下文获取用户 ID
// 兼容 uint、float64、int 等多种类型
func GetUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	// 类型断言，兼容多种可能的类型
	switch v := userID.(type) {
	case uint:
		return v
	case float64:
		return uint(v)
	case int:
		return uint(v)
	case uint32:
		return uint(v)
	case int64:
		return uint(v)
	default:
		logger.Errorf("GetUserIDFromContext: invalid user_id type: %T", userID)
		return 0
	}
}

// GetUsernameFromContext 从上下文获取用户名
func GetUsernameFromContext(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}

	if name, ok := username.(string); ok {
		return name
	}

	return ""
}
