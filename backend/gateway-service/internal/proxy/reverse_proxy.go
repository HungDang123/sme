package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewReverseProxy(rawTarget string) gin.HandlerFunc {
	target, err := url.Parse(rawTarget)
	if err != nil {
		panic(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api")
	}

	return func(c *gin.Context) {
		reverseProxy.ServeHTTP(c.Writer, c.Request)
	}
}
