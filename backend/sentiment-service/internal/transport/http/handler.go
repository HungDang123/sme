package http

import (
	nethttp "net/http"

	"sme-social-listening/sentiment-service/internal/config"
	"sme-social-listening/sentiment-service/internal/domain"
	"sme-social-listening/sentiment-service/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	analyzer *service.Analyzer
}

func NewRouter(cfg config.Config) *gin.Engine {
	handler := &Handler{analyzer: service.NewAnalyzer(cfg.GeminiAPIKey)}

	router := gin.Default()
	router.GET("/health", handler.health)
	router.POST("/sentiment/analyze", handler.analyze)
	return router
}

func (h *Handler) health(c *gin.Context) {
	ok(c, gin.H{"service": "sentiment-service", "status": "ok"})
}

func (h *Handler) analyze(c *gin.Context) {
	var req domain.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, nethttp.StatusBadRequest, "invalid request body")
		return
	}

	ok(c, h.analyzer.Analyze(req))
}
