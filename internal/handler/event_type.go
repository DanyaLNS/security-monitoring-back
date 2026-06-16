package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"security-monitor/internal/dto"
	"security-monitor/internal/service"
)

type EventTypeHandler struct {
	service *service.EventTypeService
}

func NewEventTypeHandler(
	service *service.EventTypeService,
) *EventTypeHandler {
	return &EventTypeHandler{
		service: service,
	}
}

func (h *EventTypeHandler) Create(c *gin.Context) {
	var input dto.CreateEventType

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	eventType, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, eventType)
}

func (h *EventTypeHandler) GetAll(c *gin.Context) {
	eventTypes, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, eventTypes)
}

func (h *EventTypeHandler) GetByULID(c *gin.Context) {
	ulid := c.Param("ulid")

	eventType, err := h.service.GetByULID(c.Request.Context(), ulid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, eventType)
}

func (h *EventTypeHandler) Delete(c *gin.Context) {
	ulid := c.Param("ulid")

	if err := h.service.Delete(c.Request.Context(), ulid); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
