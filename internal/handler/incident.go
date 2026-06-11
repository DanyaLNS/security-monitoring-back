package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"security-monitor/internal/service"
)

type IncidentHandler struct {
	service *service.IncidentService
}

func NewIncidentHandler(
	service *service.IncidentService,
) *IncidentHandler {
	return &IncidentHandler{
		service: service,
	}
}

func (h *IncidentHandler) GetAll(c *gin.Context) {
	result, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *IncidentHandler) GetByULID(c *gin.Context) {
	ulid := c.Param("ulid")

	result, err := h.service.GetByULID(c.Request.Context(), ulid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *IncidentHandler) GetIncidentEvents(c *gin.Context) {
	ulid := c.Param("ulid")

	result, err := h.service.GetIncidentEvents(c.Request.Context(), ulid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
