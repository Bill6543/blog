package service

import (
	"blog/pkg/errors"
	"blog/pkg/logger"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadService 文件上传服务
type UploadService struct {
	basePath string // 基础路径，如 "./static/uploads"
	maxSize  int64  // 最大文件大小（字节）
}

// ImageType 图片类型枚举
type ImageType string

const (
	ImageTypeAvatar  ImageType = "avatars"  // 头像
	ImageTypeCover   ImageType = "covers"   // 封面
	ImageTypeContent ImageType = "contents" // 内容图片
)

// UploadConfig 上传配置
type UploadConfig struct {
	Dir          ImageType                          // 上传目录类型
	MaxSize      int64                              // 最大大小（可选，覆盖默认值）
	AllowTypes   []string                           // 允许的文件类型（可选，覆盖默认值）
	GenerateName func(*multipart.FileHeader) string // 自定义文件名生成函数
}

// NewUploadService 创建上传服务实例
func NewUploadService(basePath string, maxSize int64) *UploadService {
	return &UploadService{
		basePath: basePath,
		maxSize:  maxSize,
	}
}

// UploadImage 上传图片
func (s *UploadService) UploadImage(file *multipart.FileHeader, config UploadConfig) (string, error) {
	// 1. 验证文件类型
	if err := s.validateFileType(file, config); err != nil {
		logger.Warnf("File upload failed: invalid file type, filename=%s, error=%v", file.Filename, err)
		return "", errors.New(errors.InvalidFileTypeCode)
	}

	// 2. 验证文件大小
	if err := s.validateFileSize(file, config); err != nil {
		logger.Warnf("File upload failed: file too large, filename=%s, size=%d bytes, error=%v", file.Filename, file.Size, err)
		return "", errors.New(errors.FileTooLargeCode)
	}

	// 3. 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", errors.New(errors.FileOpenFailedCode)
	}
	defer src.Close()

	// 4. 生成文件名
	filename := s.generateFilename(file, config)

	// 5. 构建完整路径
	uploadDir := filepath.Join(s.basePath, string(config.Dir))
	uploadPath := filepath.Join(uploadDir, filename)

	// 6. 确保目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Errorf("File upload failed: create directory failed, path=%s, error=%v", uploadDir, err)
		return "", errors.New(errors.FileSaveFailedCode)
	}

	// 7. 创建目标文件
	dst, err := os.Create(uploadPath)
	if err != nil {
		logger.Errorf("File upload failed: create file failed, path=%s, error=%v", uploadPath, err)
		return "", errors.New(errors.FileSaveFailedCode)
	}
	defer dst.Close()

	// 8. 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		logger.Errorf("File upload failed: save file failed, path=%s, error=%v", uploadPath, err)
		return "", errors.New(errors.FileSaveFailedCode)
	}

	// 9. 返回访问 URL
	logger.Debugf("File uploaded successfully: filename=%s, path=%s, size=%d bytes", file.Filename, uploadPath, file.Size)
	return fmt.Sprintf("/static/uploads/%s/%s", config.Dir, filename), nil
}

// validateFileType 验证文件类型
func (s *UploadService) validateFileType(file *multipart.FileHeader, config UploadConfig) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// 默认允许的类型
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif"}
	if len(config.AllowTypes) > 0 {
		allowedTypes = config.AllowTypes
	}

	for _, t := range allowedTypes {
		if ext == t {
			return nil
		}
	}

	return fmt.Errorf("不支持的图片格式，仅支持：%s", strings.Join(allowedTypes, ", "))
}

// validateFileSize 验证文件大小
func (s *UploadService) validateFileSize(file *multipart.FileHeader, config UploadConfig) error {
	maxSize := s.maxSize
	if config.MaxSize > 0 {
		maxSize = config.MaxSize
	}

	if file.Size > maxSize {
		return fmt.Errorf("文件过大，最大允许 %dMB", maxSize/1024/1024)
	}

	return nil
}

// generateFilename 生成文件名
func (s *UploadService) generateFilename(file *multipart.FileHeader, config UploadConfig) string {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// 如果提供了自定义生成函数，使用它
	if config.GenerateName != nil {
		return config.GenerateName(file) + ext
	}

	// 默认使用时间戳 + 随机数
	filename := fmt.Sprintf("%s_%d%s",
		time.Now().Format("20060102150405"),
		time.Now().UnixNano()%10000,
		ext,
	)

	return filename
}

// DeleteImage 删除图片
func (s *UploadService) DeleteImage(imageURL string) error {
	// 将 URL 转换为文件路径
	filePath := "." + imageURL

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在，直接返回
	}

	// 删除文件
	return os.Remove(filePath)
}

// GetDefaultAvatar 获取默认头像
func (s *UploadService) GetDefaultAvatar() string {
	return "/static/default_avatar.png"
}
