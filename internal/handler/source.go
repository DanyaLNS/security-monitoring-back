package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"security-monitor/internal/domain"
	"security-monitor/internal/service"
)

type SourceHandler struct {
	service *service.SourceService
}

func NewSourceHandler(
	service *service.SourceService,
) *SourceHandler {
	return &SourceHandler{
		service: service,
	}
}

func (h *SourceHandler) Create(c *gin.Context) {
	var input domain.EventSource

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SourceHandler) GetAll(c *gin.Context) {
	result, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SourceHandler) GetByULID(c *gin.Context) {
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

func (h *SourceHandler) Delete(c *gin.Context) {
	ulid := c.Param("ulid")

	if err := h.service.Delete(c.Request.Context(), ulid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
