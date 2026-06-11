package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"security-monitor/internal/service"
)

type AnalysisHandler struct {
	service *service.AnalysisService
}

func NewAnalysisHandler(
	service *service.AnalysisService,
) *AnalysisHandler {
	return &AnalysisHandler{
		service: service,
	}
}

func (h *AnalysisHandler) GetAll(c *gin.Context) {
	result, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AnalysisHandler) GetByEventULID(c *gin.Context) {
	eventULID := c.Param("event_ulid")

	result, err := h.service.GetByEventULID(
		c.Request.Context(),
		eventULID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
