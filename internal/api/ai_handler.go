package api

import (
	"blog/pkg/logger"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// AIHandler AI Handler
type AIHandler struct {
	*Handler
}

// GenerateSummaryRequest 生成摘要请求
type GenerateSummaryRequest struct {
	Content string `json:"content" binding:"required"`
}

// GenerateSummaryResponse 生成摘要响应
type GenerateSummaryResponse struct {
	Summary string `json:"summary"`
}

// GenerateCoverRequest 生成封面请求
type GenerateCoverRequest struct {
	Content string `json:"content" binding:"required"`
}

// GenerateCoverResponse 生成封面响应
type GenerateCoverResponse struct {
	CoverURL string `json:"cover_url"`
}

// GenerateSummary 生成文章摘要
func (h *AIHandler) GenerateSummary(c *gin.Context) {
	var req GenerateSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败："+err.Error())
		return
	}

	summary, err := h.AIService.GenerateSummary(req.Content)
	if err != nil {
		logger.Errorf("Generate summary failed: error=%v", err)
		response.InternalError(c, "生成摘要失败")
		return
	}

	resp := GenerateSummaryResponse{
		Summary: summary,
	}

	logger.Debugf("Generate summary request received")
	response.Success(c, resp)
}

// GenerateCover 生成文章封面
func (h *AIHandler) GenerateCover(c *gin.Context) {
	var req GenerateCoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败："+err.Error())
		return
	}

	coverURL, err := h.AIService.GenerateCover(req.Content)
	if err != nil {
		logger.Errorf("Generate cover failed: error=%v", err)
		response.InternalError(c, "生成封面失败")
		return
	}

	resp := GenerateCoverResponse{
		CoverURL: coverURL,
	}

	logger.Debugf("Generate cover request received")
	response.Success(c, resp)
}
