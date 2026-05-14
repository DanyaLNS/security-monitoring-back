package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"security-monitor/internal/dto"
	"security-monitor/internal/service"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(
	service *service.EventService,
) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

func (h *EventHandler) Create(c *gin.Context) {
	var req dto.CreateEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *EventHandler) GetAll(c *gin.Context) {
	events, err := h.service.GetAll(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) Delete(c *gin.Context) {
	ulid := c.Param("ulid")

	err := h.service.Delete(
		c.Request.Context(),
		ulid,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
