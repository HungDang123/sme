package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
)

func ok(c *gin.Context, data any) {
	c.JSON(nethttp.StatusOK, gin.H{"data": data})
}

func fail(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
