package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"security-monitor/internal/service"
)

type DashboardHandler struct {
	service *service.DashboardService
}

func NewDashboardHandler(
	service *service.DashboardService,
) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

func (h *DashboardHandler) GetMetrics(c *gin.Context) {
	result, err := h.service.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DashboardHandler) GetTimeline(c *gin.Context) {
	result, err := h.service.GetTimeline(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DashboardHandler) GetSeverityDistribution(c *gin.Context) {
	result, err := h.service.GetSeverityDistribution(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DashboardHandler) GetTopSources(c *gin.Context) {
	result, err := h.service.GetTopSources(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DashboardHandler) GetRecentEvents(c *gin.Context) {
	result, err := h.service.GetRecentEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
