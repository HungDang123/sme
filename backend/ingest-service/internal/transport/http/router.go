package http

import (
	nethttp "net/http"

	"sme-social-listening/ingest-service/internal/crawler"

	"github.com/gin-gonic/gin"
)

func NewRouter(c *crawler.Crawler) *gin.Engine {
	router := gin.Default()
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(nethttp.StatusOK, gin.H{"data": gin.H{"service": "ingest-service", "status": "ok"}})
	})
	router.POST("/ingest/trigger", func(ctx *gin.Context) {
		go func() {
			if err := c.Run(ctx.Request.Context()); err != nil {
				// logged in crawler
			}
		}()
		ctx.JSON(nethttp.StatusAccepted, gin.H{"data": gin.H{"status": "crawl started"}})
	})
	return router
}
