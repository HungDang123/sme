package http

import (
	nethttp "net/http"

	"sme-social-listening/gateway-service/internal/config"
	"sme-social-listening/gateway-service/internal/proxy"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(cors())

	router.GET("/health", func(c *gin.Context) {
		ok(c, gin.H{"service": "gateway-service", "status": "ok"})
	})

	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		ok(c, gin.H{"service": "gateway-service", "status": "ok"})
	})

	mountProxy(api, "/brands", cfg.BrandServiceURL)
	mountProxy(api, "/mentions", cfg.MentionServiceURL)
	mountProxy(api, "/sentiment", cfg.SentimentServiceURL)
	mountProxy(api, "/alerts", cfg.AlertServiceURL)
	mountProxy(api, "/ingest", cfg.IngestServiceURL)

	return router
}

func mountProxy(group *gin.RouterGroup, prefix string, target string) {
	handler := proxy.NewReverseProxy(target)
	group.Any(prefix, handler)
	group.Any(prefix+"/*path", handler)
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization,X-Ingest-Token")
		if c.Request.Method == nethttp.MethodOptions {
			c.AbortWithStatus(nethttp.StatusNoContent)
			return
		}
		c.Next()
	}
}

func ok(c *gin.Context, data any) {
	c.JSON(nethttp.StatusOK, gin.H{"data": data})
}
