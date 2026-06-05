package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LogConfig 日志配置
type LogConfig struct {
	// 是否记录请求体
	EnableRequestBody bool `yaml:"enable_request_body"`
	// 是否记录响应体
	EnableResponseBody bool `yaml:"enable_response_body"`
	// 慢查询阈值（毫秒）
	SlowRequestThreshold int `yaml:"slow_request_threshold"`
	// 需要脱敏的字段
	SensitiveFields []string `yaml:"sensitive_fields"`
}

// responseWriter 自定义响应写入器，用于捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 增强版日志中间件
func Logger(config LogConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求 ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		// 记录开始时间
		start := time.Now()

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()

		// 记录请求体（如果启用）
		var requestBody string
		if config.EnableRequestBody && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// 脱敏处理
				requestBody = sensitiveFilter(string(bodyBytes), config.SensitiveFields)
				// 重新写入请求体，避免影响后续处理
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 创建自定义响应写入器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// 执行请求
		c.Next()

		// 计算请求耗时
		latency := time.Since(start)
		latencyMs := latency.Milliseconds()

		// 获取状态码
		statusCode := c.Writer.Status()

		// 获取错误信息
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// 记录响应体（如果启用）
		var responseBody string
		if config.EnableResponseBody {
			responseBody = sensitiveFilter(writer.body.String(), config.SensitiveFields)
		}

		// 检测慢查询
		isSlow := latencyMs > int64(config.SlowRequestThreshold)

		// 构建日志字段
		logData := map[string]interface{}{
			"timestamp":  time.Now().Format(time.RFC3339),
			"request_id": requestID,
			"method":     method,
			"path":       path,
			"query":      query,
			"status":     statusCode,
			"latency_ms": latencyMs,
			"client_ip":  clientIP,
			"user_agent": userAgent,
			"referer":    referer,
			"error":      errorMessage,
		}

		// 添加请求体（如果启用）
		if config.EnableRequestBody {
			logData["request_body"] = requestBody
		}

		// 添加响应体（如果启用）
		if config.EnableResponseBody {
			logData["response_body"] = responseBody
		}

		// 添加慢查询标记
		if isSlow {
			logData["is_slow"] = true
		}

		// 根据状态码选择日志级别
		if statusCode >= http.StatusInternalServerError {
			logger.Errorf("[Server Error] %+v", logData)
		} else if statusCode >= http.StatusBadRequest {
			logger.Warnf("[Client Error] %+v", logData)
		} else if isSlow {
			logger.Warnf("[Slow Request] %+v", logData)
		} else {
			logger.Infof("[Request] %+v", logData)
		}
	}
}

// sensitiveFilter 敏感信息过滤
func sensitiveFilter(data string, sensitiveFields []string) string {
	// 尝试解析为 JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		// 不是 JSON 格式，尝试简单的字符串替换
		result := data
		for _, field := range sensitiveFields {
			// 简单的键值对替换（针对 form-data 等格式）
			patterns := []string{
				field + "=",
				field + ":",
				"\"" + field + "\"",
			}
			for _, pattern := range patterns {
				if strings.Contains(result, pattern) {
					// 检测到敏感信息，返回脱敏标记
					return "***REDACTED***"
				}
			}
		}
		return data
	}

	// JSON 格式脱敏处理
	for _, field := range sensitiveFields {
		if _, exists := jsonData[field]; exists {
			jsonData[field] = "***REDACTED***"
		}
	}

	result, err := json.Marshal(jsonData)
	if err != nil {
		return data
	}

	return string(result)
}

// DefaultLogConfig 默认日志配置
func DefaultLogConfig() LogConfig {
	return LogConfig{
		EnableRequestBody:    false, // 默认不记录请求体（性能考虑）
		EnableResponseBody:   false, // 默认不记录响应体
		SlowRequestThreshold: 1000,  // 1 秒阈值
		SensitiveFields:      []string{"password", "token", "secret", "authorization"},
	}
}
