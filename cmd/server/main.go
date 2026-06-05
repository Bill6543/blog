package main

import (
	"blog/internal/api"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/internal/router"
	"blog/internal/service"
	"blog/pkg/ai"
	"blog/pkg/config"
	"blog/pkg/database"
	"blog/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志系统（必须在最前面）
	if err := logger.InitLogger(
		cfg.Log.Level,
		cfg.Log.FilePath,
		cfg.Log.MaxSize,
		cfg.Log.MaxBackups,
		cfg.Log.MaxAge,
	); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close() // 程序退出时关闭日志

	// 使用 zap 记录启动日志
	logger.Infof("Starting application with config: %+v", cfg)

	// 3. 初始化数据库
	if err := database.InitDB(&cfg.Database); err != nil {
		logger.Errorf("Failed to init database: %v", err)
		logger.Sync() // 确保刷新日志
		os.Exit(1)
	}

	// 4. 自动迁移模型
	if err := database.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Tag{},
		&entity.Article{},
		&entity.ArticleTag{},
		&entity.Comment{},
		&entity.Like{},
	); err != nil {
		logger.Errorf("Failed to auto migrate: %v", err)
		logger.Sync() // 确保刷新日志
		os.Exit(1)
	}

	// 5. 初始化 Repository
	userRepo := repository.NewUserRepository(database.GetDB())
	categoryRepo := repository.NewCategoryRepository(database.GetDB())
	tagRepo := repository.NewTagRepository(database.GetDB())
	articleRepo := repository.NewArticleRepository(database.GetDB())
	commentRepo := repository.NewCommentRepository(database.GetDB())
	likeRepo := repository.NewLikeRepository(database.GetDB())

	// 6. 初始化 Service
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpireTime)
	userService := service.NewUserService(userRepo)
	uploadService := service.NewUploadService("./static/uploads", 2*1024*1024)
	articleService := service.NewArticleService(articleRepo, tagRepo, categoryRepo, database.GetDB())
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	commentService := service.NewCommentService(commentRepo, articleRepo, database.GetDB())
	likeService := service.NewLikeService(likeRepo, articleRepo, userRepo, database.GetDB())
	aiService := ai.NewCozeService(cfg.Coze.APIKey, cfg.Coze.BotID, cfg.Coze.APIURL)

	// 7. 初始化 Handler
	handler := api.NewHandler(authService, userService, articleService, categoryService, tagService, commentService, likeService, uploadService, aiService)

	// 8. 设置路由
	r := router.SetupRouter(handler, cfg, userRepo)

	// 9. 创建 HTTP 服务器实例
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在 goroutine 中启动服务
	go func() {
		logger.Infof("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Failed to start server: %v", err)
		}
	}()

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭 HTTP 服务器（给 5 秒宽限期）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server stopped")
}
