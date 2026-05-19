package http

import (
	"errors"
	nethttp "net/http"
	"strings"

	"sme-social-listening/brand-service/internal/service"

	"github.com/gin-gonic/gin"
)

type createBrandRequest struct {
	Name             string   `json:"name"`
	Keywords         []string `json:"keywords"`
	TelegramChatID   string   `json:"telegramChatId"`
}

type addKeywordRequest struct {
	Value string `json:"value"`
}

type updateTelegramRequest struct {
	TelegramChatID string `json:"telegramChatId"`
}

type Handler struct {
	brandService *service.BrandService
}

func NewHandler(brandService *service.BrandService) *Handler {
	return &Handler{brandService: brandService}
}

func NewRouter(brandService *service.BrandService) *gin.Engine {
	handler := NewHandler(brandService)
	router := gin.Default()
	router.GET("/health", handler.health)
	router.GET("/brands", handler.listBrands)
	router.GET("/brands/:id", handler.getBrand)
	router.POST("/brands", handler.createBrand)
	router.GET("/brands/:id/keywords", handler.listKeywords)
	router.POST("/brands/:id/keywords", handler.addKeyword)
	router.PATCH("/brands/:id/telegram", handler.updateTelegram)
	router.GET("/keywords/active", handler.listActiveKeywords)
	return router
}

func (h *Handler) health(c *gin.Context) {
	ok(c, gin.H{"service": "brand-service", "status": "ok"})
}

func (h *Handler) listBrands(c *gin.Context) {
	brands, err := h.brandService.List(c.Request.Context())
	if err != nil {
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}
	ok(c, brands)
}

func (h *Handler) getBrand(c *gin.Context) {
	brand, err := h.brandService.GetBrand(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrBrandNotFound) {
			fail(c, nethttp.StatusNotFound, "brand not found")
			return
		}
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}
	ok(c, brand)
}

func (h *Handler) createBrand(c *gin.Context) {
	var req createBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, nethttp.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		fail(c, nethttp.StatusBadRequest, "brand name is required")
		return
	}

	brand, err := h.brandService.Create(c.Request.Context(), req.Name, req.Keywords)
	if err != nil {
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}
	if req.TelegramChatID != "" {
		_ = h.brandService.UpdateTelegramChatID(c.Request.Context(), brand.ID, req.TelegramChatID)
		brand, _ = h.brandService.GetBrand(c.Request.Context(), brand.ID)
	}
	c.JSON(nethttp.StatusCreated, gin.H{"data": brand})
}

func (h *Handler) listKeywords(c *gin.Context) {
	keywords, err := h.brandService.ListKeywords(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrBrandNotFound) {
			fail(c, nethttp.StatusNotFound, "brand not found")
			return
		}
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}
	ok(c, keywords)
}

func (h *Handler) addKeyword(c *gin.Context) {
	var req addKeywordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, nethttp.StatusBadRequest, "invalid request body")
		return
	}
	keyword, err := h.brandService.AddKeyword(c.Request.Context(), c.Param("id"), req.Value)
	if err != nil {
		if errors.Is(err, service.ErrBrandNotFound) {
			fail(c, nethttp.StatusNotFound, "brand not found")
			return
		}
		fail(c, nethttp.StatusBadRequest, err.Error())
		return
	}
	c.JSON(nethttp.StatusCreated, gin.H{"data": keyword})
}

func (h *Handler) updateTelegram(c *gin.Context) {
	var req updateTelegramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, nethttp.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.brandService.UpdateTelegramChatID(c.Request.Context(), c.Param("id"), req.TelegramChatID); err != nil {
		if errors.Is(err, service.ErrBrandNotFound) {
			fail(c, nethttp.StatusNotFound, "brand not found")
			return
		}
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}
	brand, _ := h.brandService.GetBrand(c.Request.Context(), c.Param("id"))
	ok(c, brand)
}

func (h *Handler) listActiveKeywords(c *gin.Context) {
	keywords, err := h.brandService.ListActiveKeywords(c.Request.Context())
	if err != nil {
		fail(c, nethttp.StatusInternalServerError, err.Error())
		return
	}
	ok(c, keywords)
}
