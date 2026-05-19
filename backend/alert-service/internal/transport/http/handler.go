package http

import (
	nethttp "net/http"

	"sme-social-listening/alert-service/internal/domain"
	"sme-social-listening/alert-service/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	alertService *service.AlertService
}

func NewRouter(alertService *service.AlertService) *gin.Engine {
	handler := &Handler{alertService: alertService}

	router := gin.Default()
	router.GET("/health", handler.health)
	router.GET("/alerts", handler.listAlerts)
	router.POST("/alerts", handler.createAlert)
	return router
}

func (h *Handler) health(c *gin.Context) {
	ok(c, gin.H{"service": "alert-service", "status": "ok"})
}

func (h *Handler) listAlerts(c *gin.Context) {
	alerts, err := h.alertService.List(c.Request.Context())
	if err != nil {
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}
	ok(c, alerts)
}

func (h *Handler) createAlert(c *gin.Context) {
	var req domain.CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, nethttp.StatusBadRequest, "invalid request body")
		return
	}

	alert, err := h.alertService.Create(c.Request.Context(), req)
	if err != nil {
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(nethttp.StatusCreated, gin.H{"data": alert})
}
