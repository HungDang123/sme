package http

import (
	nethttp "net/http"
	"strconv"
	"strings"
	"time"

	"sme-social-listening/mention-service/internal/config"
	"sme-social-listening/mention-service/internal/domain"
	"sme-social-listening/mention-service/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	mentionService *service.MentionService
	ingestToken    string
}

func NewRouter(mentionService *service.MentionService, cfg config.Config) *gin.Engine {
	handler := &Handler{mentionService: mentionService, ingestToken: cfg.IngestToken}
	router := gin.Default()
	router.GET("/health", handler.health)
	router.GET("/mentions", handler.listMentions)
	router.POST("/mentions", handler.createMention)
	return router
}

func (h *Handler) health(c *gin.Context) {
	ok(c, gin.H{"service": "mention-service", "status": "ok"})
}

func (h *Handler) listMentions(c *gin.Context) {
	filter := domain.MentionFilter{
		BrandID:   c.Query("brandId"),
		Keyword:   c.Query("keyword"),
		Source:    c.Query("source"),
		Sentiment: c.Query("sentiment"),
	}
	if from := c.Query("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			filter.From = &t
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			filter.To = &t
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	mentions, err := h.mentionService.List(c.Request.Context(), filter)
	if err != nil {
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}
	ok(c, mentions)
}

func (h *Handler) createMention(c *gin.Context) {
	var req domain.CreateMentionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, nethttp.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		fail(c, nethttp.StatusBadRequest, "content is required")
		return
	}

	source := strings.TrimSpace(req.Source)
	if source == "" {
		source = "manual"
		req.Source = source
	}
	if source != "manual" && h.ingestToken != "" {
		if c.GetHeader("X-Ingest-Token") != h.ingestToken {
			fail(c, nethttp.StatusUnauthorized, "invalid ingest token")
			return
		}
	}

	result, err := h.mentionService.Create(c.Request.Context(), req)
	if err != nil {
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}

	status := nethttp.StatusCreated
	if !result.Created {
		status = nethttp.StatusOK
	}
	c.JSON(status, gin.H{"data": result.Mention, "created": result.Created})
}
