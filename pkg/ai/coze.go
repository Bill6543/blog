package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"blog/pkg/logger"
)

// 未配置超时时的兜底值
const (
	defaultTimeout      = 30 * time.Second
	defaultCoverTimeout = 120 * time.Second
)

// CozeConfig Coze配置
type CozeConfig struct {
	APIKey string
	BotID  string
	APIURL string
	// Timeout 摘要生成超时；CoverTimeout 封面生成超时（文生图耗时更长，单独配置）
	Timeout      time.Duration
	CoverTimeout time.Duration
}

// CozeService Coze AI服务
type CozeService struct {
	config *CozeConfig
	client *http.Client
}

// NewCozeService 创建Coze服务实例
// timeout/coverTimeout 由配置传入，传 0 时使用默认值
func NewCozeService(apiKey, botID, apiURL string, timeout, coverTimeout time.Duration) *CozeService {
	if apiURL == "" {
		apiURL = "https://api.coze.cn/open_api/v2/chat"
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	if coverTimeout <= 0 {
		coverTimeout = defaultCoverTimeout
	}

	config := &CozeConfig{
		APIKey:       apiKey,
		BotID:        botID,
		APIURL:       apiURL,
		Timeout:      timeout,
		CoverTimeout: coverTimeout,
	}

	return &CozeService{
		config: config,
		// 超时由每次请求的 context 控制，两个接口取值不同，故不设 client 级 Timeout
		client: &http.Client{},
	}
}

// ChatRequest Coze聊天请求
type ChatRequest struct {
	BotID      string `json:"bot_id"`
	User       string `json:"user"`
	Query      string `json:"query"`
	Stream     bool   `json:"stream"`
	AutoSave   bool   `json:"auto_save"`
	Additional []struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"additional_messages,omitempty"`
}

// ChatResponse Coze聊天响应
type ChatResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"msg"`
	Messages []struct {
		Role        string `json:"role"`
		Type        string `json:"type"`
		Content     string `json:"content"`
		ContentType string `json:"content_type,omitempty"`
	} `json:"messages"`
	ConversationID string `json:"conversation_id"`
}

// GenerateSummary 生成文章摘要
func (s *CozeService) GenerateSummary(content string) (string, error) {
	if s.config.APIKey == "" || s.config.BotID == "" {
		logger.Warn("Coze API key or Bot ID is not configured")
		return "", fmt.Errorf("Coze API key or Bot ID is not configured")
	}

	prompt := fmt.Sprintf("请为以下文章生成一个简洁的摘要（不超过200字）：\n\n%s", content)

	req := ChatRequest{
		BotID:    s.config.BotID,
		User:     "blog_system",
		Query:    prompt,
		Stream:   false,
		AutoSave: true,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("Failed to marshal Coze request: %v", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	logger.Debugf("Coze API request: bot_id=%s, body=%s", s.config.BotID, string(reqBody))

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.config.APIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Errorf("Failed to create HTTP request: %v", err)
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.APIKey))

	resp, err := s.client.Do(httpReq)
	if err != nil {
		logger.Errorf("Failed to send Coze request: %v", err)
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to read Coze response: %v", err)
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	logger.Debugf("Coze API response: status=%d, body=%s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Coze API returned error status: %d, body: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("Coze API returned error status: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		logger.Errorf("Failed to unmarshal Coze response: %v", err)
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chatResp.Code != 0 {
		logger.Errorf("Coze API returned error code: %d, message: %s", chatResp.Code, chatResp.Message)
		return "", fmt.Errorf("Coze API error: %s", chatResp.Message)
	}

	// 从消息数组中提取答案
	var answer string
	for _, msg := range chatResp.Messages {
		if msg.Type == "answer" && msg.Content != "" {
			answer = msg.Content
			break
		}
	}

	if answer == "" {
		logger.Errorf("No answer found in Coze response")
		return "", fmt.Errorf("no answer found in Coze response")
	}

	logger.Debugf("Coze summary generated successfully")
	return answer, nil
}

// GenerateCover 生成文章封面
func (s *CozeService) GenerateCover(content string) (string, error) {
	if s.config.APIKey == "" || s.config.BotID == "" {
		logger.Warn("Coze API key or Bot ID is not configured")
		return "", fmt.Errorf("Coze API key or Bot ID is not configured")
	}

	prompt := fmt.Sprintf("请为以下文章生成封面图片：\n\n%s", content)

	req := ChatRequest{
		BotID:    s.config.BotID,
		User:     "blog_system",
		Query:    prompt,
		Stream:   false,
		AutoSave: true,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("Failed to marshal Coze request: %v", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	logger.Debugf("Coze API request for cover: bot_id=%s, body=%s", s.config.BotID, string(reqBody))

	ctx, cancel := context.WithTimeout(context.Background(), s.config.CoverTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.config.APIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Errorf("Failed to create HTTP request: %v", err)
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.APIKey))

	resp, err := s.client.Do(httpReq)
	if err != nil {
		logger.Errorf("Failed to send Coze request: %v", err)
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to read Coze response: %v", err)
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	logger.Debugf("Coze API response for cover: status=%d, body=%s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Coze API returned error status: %d, body: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("Coze API returned error status: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		logger.Errorf("Failed to unmarshal Coze response: %v", err)
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chatResp.Code != 0 {
		logger.Errorf("Coze API returned error code: %d, message: %s", chatResp.Code, chatResp.Message)
		return "", fmt.Errorf("Coze API error: %s", chatResp.Message)
	}

	// 从消息数组中提取图片URL
	// Coze返回的图片URL在 tool_response 类型的消息中，content字段包含JSON格式的output
	var coverURL string
	for _, msg := range chatResp.Messages {
		if msg.Type == "tool_response" && msg.Content != "" {
			// 解析tool_response中的JSON，提取output字段
			var toolResp struct {
				Output string `json:"output"`
			}
			if err := json.Unmarshal([]byte(msg.Content), &toolResp); err == nil {
				if toolResp.Output != "" {
					coverURL = toolResp.Output
					break
				}
			}
		}
	}

	// 如果没有找到tool_response，尝试从answer中提取URL
	if coverURL == "" {
		for _, msg := range chatResp.Messages {
			if msg.Type == "answer" && msg.Content != "" {
				// 尝试从answer中提取URL（可能包含描述文本）
				answer := msg.Content
				// 简单的URL提取：查找https://开头的链接
				lines := strings.Split(answer, "\n")
				for _, line := range lines {
					if strings.HasPrefix(strings.TrimSpace(line), "https://") {
						// 提取纯URL，去除可能的描述
						url := strings.TrimSpace(line)
						// 如果URL后面有其他文本，只取URL部分
						if idx := strings.Index(url, " "); idx > 0 {
							url = url[:idx]
						}
						coverURL = url
						break
					}
				}
				if coverURL != "" {
					break
				}
			}
		}
	}

	if coverURL == "" {
		logger.Errorf("No cover URL found in Coze response")
		return "", fmt.Errorf("no cover URL found in Coze response")
	}

	logger.Debugf("Coze cover generated successfully: url=%s", coverURL)
	return coverURL, nil
}
