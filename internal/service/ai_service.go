package service

import (
	"time"

	"blog/pkg/ai"
	"blog/pkg/errors"
	"blog/pkg/logger"
)

// AIService AI服务
type AIService struct {
	cozeService *ai.CozeService
}

// NewAIService 创建AI服务实例
func NewAIService(apiKey, botID, apiURL string, timeout, coverTimeout time.Duration) *AIService {
	cozeService := ai.NewCozeService(apiKey, botID, apiURL, timeout, coverTimeout)
	return &AIService{
		cozeService: cozeService,
	}
}

// GenerateSummary 生成文章摘要
func (s *AIService) GenerateSummary(content string) (string, error) {
	if content == "" {
		logger.Warn("Generate summary failed: content is empty")
		return "", errors.New(errors.InvalidParam)
	}

	summary, err := s.cozeService.GenerateSummary(content)
	if err != nil {
		logger.Errorf("Generate summary failed: error=%v", err)
		return "", err
	}

	logger.Debugf("Summary generated successfully: length=%d", len(summary))
	return summary, nil
}

// GenerateCover 生成文章封面
func (s *AIService) GenerateCover(content string) (string, error) {
	if content == "" {
		logger.Warn("Generate cover failed: content is empty")
		return "", errors.New(errors.InvalidParam)
	}

	coverURL, err := s.cozeService.GenerateCover(content)
	if err != nil {
		logger.Errorf("Generate cover failed: error=%v", err)
		return "", err
	}

	logger.Debugf("Cover generated successfully: url=%s", coverURL)
	return coverURL, nil
}
