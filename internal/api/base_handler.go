package api

import (
	"blog/internal/service"
	"blog/pkg/ai"
)

// Handler 所有 Handler 的基类
type Handler struct {
	AuthService     *service.AuthService
	UserService     *service.UserService
	ArticleService  *service.ArticleService
	CategoryService *service.CategoryService
	TagService      *service.TagService
	CommentService  *service.CommentService
	LikeService     *service.LikeService
	UploadService   *service.UploadService
	AIService       *ai.CozeService
}

// NewHandler 创建 Handler 实例
func NewHandler(
	authService *service.AuthService,
	userService *service.UserService,
	articleService *service.ArticleService,
	categoryService *service.CategoryService,
	tagService *service.TagService,
	commentService *service.CommentService,
	likeService *service.LikeService,
	uploadService *service.UploadService,
	aiService *ai.CozeService,
) *Handler {
	return &Handler{
		AuthService:     authService,
		UserService:     userService,
		ArticleService:  articleService,
		CategoryService: categoryService,
		TagService:      tagService,
		CommentService:  commentService,
		LikeService:     likeService,
		UploadService:   uploadService,
		AIService:       aiService,
	}
}
