package router

import (
	"blog/internal/api"
	"blog/internal/middleware"
	"blog/internal/repository"
	"blog/pkg/config"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter(handler *api.Handler, cfg *config.AppConfig, userRepo *repository.UserRepository) *gin.Engine {
	r := gin.New()

	// 恢复中间件（使用默认配置）
	r.Use(gin.Recovery())

	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(middleware.LogConfig{
		EnableRequestBody:    true,
		EnableResponseBody:   true,
		SlowRequestThreshold: 1000,
		SensitiveFields:      []string{"password", "token"},
	}))

	// 静态文件
	r.Static("/static", "./static")

	// API 路由组
	apiGroup := r.Group("/api")
	{
		// 认证路由（无需登录）
		auth := apiGroup.Group("/auth")
		{
			auth.POST("/register", (&api.AuthHandler{Handler: handler}).Register)
			auth.POST("/login", (&api.AuthHandler{Handler: handler}).Login)
		}

		// 公开路由（无需登录）
		public := apiGroup.Group("")
		{
			// 文章
			public.GET("/articles", (&api.ArticleHandler{Handler: handler}).GetArticleList)
			public.GET("/articles/:id", (&api.ArticleHandler{Handler: handler}).GetArticle)

			// 分类
			public.GET("/categories", (&api.CategoryHandler{Handler: handler}).GetCategoryList)
			public.GET("/categories/:id", (&api.CategoryHandler{Handler: handler}).GetCategory)

			// 标签
			public.GET("/tags", (&api.TagHandler{Handler: handler}).GetTagList)
			public.GET("/tags/:id", (&api.TagHandler{Handler: handler}).GetTag)

			// 评论
			public.GET("/comments", (&api.CommentHandler{Handler: handler}).GetCommentList)
			public.GET("/comments/:id", (&api.CommentHandler{Handler: handler}).GetComment)
		}

		// 需要认证的路由
		protected := apiGroup.Group("")
		protected.Use(middleware.Auth(middleware.AuthConfig{
			Secret:   cfg.JWT.Secret,
			UserRepo: userRepo,
		}))
		{
			// 认证相关
			protected.GET("/auth/me", (&api.AuthHandler{Handler: handler}).GetCurrentUser)
			protected.POST("/auth/logout", (&api.AuthHandler{Handler: handler}).Logout)

			// 查看指定用户的文章
			protected.GET("/users/:id/articles", (&api.ArticleHandler{Handler: handler}).GetUserArticles)

			// 用户管理（智能权限判断）
			protected.PUT("/users/:id", (&api.UserHandler{Handler: handler}).UpdateUserInfo)

			// 文章管理（普通用户 + 管理员，需校验资源所有权）
			protected.POST("/articles", (&api.ArticleHandler{Handler: handler}).CreateArticle)
			protected.PUT("/articles/:id", (&api.ArticleHandler{Handler: handler}).UpdateArticle)
			protected.DELETE("/articles/:id", (&api.ArticleHandler{Handler: handler}).DeleteArticle)
			protected.POST("/articles/:id/restore", (&api.ArticleHandler{Handler: handler}).RestoreArticle)
			protected.POST("/articles/:id/view", (&api.ArticleHandler{Handler: handler}).IncreaseViewCount)

			// 点赞管理（需要登录）
			protected.POST("/articles/:id/like", (&api.LikeHandler{Handler: handler}).LikeArticle)
			protected.DELETE("/articles/:id/like", (&api.LikeHandler{Handler: handler}).UnlikeArticle)
			protected.GET("/articles/:id/like", (&api.LikeHandler{Handler: handler}).GetLikeStatus)

			// 评论管理（普通用户 + 管理员，需校验资源所有权）
			protected.POST("/comments", (&api.CommentHandler{Handler: handler}).CreateComment)
			protected.PUT("/comments/:id", (&api.CommentHandler{Handler: handler}).UpdateComment)
			protected.DELETE("/comments/:id", (&api.CommentHandler{Handler: handler}).DeleteComment)

			// AI 服务（需要登录）
			protected.POST("/ai/summary", (&api.AIHandler{Handler: handler}).GenerateSummary)
			protected.POST("/ai/cover", (&api.AIHandler{Handler: handler}).GenerateCover)
		}

		// 管理员专属路由
		adminOnly := apiGroup.Group("")
		adminOnly.Use(middleware.Auth(middleware.AuthConfig{
			Secret:   cfg.JWT.Secret,
			UserRepo: userRepo,
		}))
		adminOnly.Use(middleware.IsAdmin())
		{
			// 用户管理（仅管理员）
			adminOnly.GET("/users", (&api.UserHandler{Handler: handler}).GetUserList)
			adminOnly.DELETE("/users/:id", (&api.UserHandler{Handler: handler}).DeleteUser)

			// 分类管理（仅管理员）
			adminOnly.POST("/categories", (&api.CategoryHandler{Handler: handler}).CreateCategory)
			adminOnly.PUT("/categories/:id", (&api.CategoryHandler{Handler: handler}).UpdateCategory)
			adminOnly.DELETE("/categories/:id", (&api.CategoryHandler{Handler: handler}).DeleteCategory)

			// 标签管理（仅管理员）
			adminOnly.POST("/tags", (&api.TagHandler{Handler: handler}).CreateTag)
			adminOnly.PUT("/tags/:id", (&api.TagHandler{Handler: handler}).UpdateTag)
			adminOnly.DELETE("/tags/:id", (&api.TagHandler{Handler: handler}).DeleteTag)

			// 评论审核（仅管理员）
			adminOnly.GET("/admin/comments/pending", (&api.CommentHandler{Handler: handler}).GetPendingComments)
			adminOnly.GET("/admin/comments/pending/count", (&api.CommentHandler{Handler: handler}).GetPendingCommentCount)
			adminOnly.PUT("/admin/comments/:id/approve", (&api.CommentHandler{Handler: handler}).ApproveComment)
			adminOnly.PUT("/admin/comments/:id/reject", (&api.CommentHandler{Handler: handler}).RejectComment)
		}
	}

	// 404 处理
	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "接口不存在")
	})

	return r
}
