package handler

import (
	"net/http"
	"strconv"

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
	filter := dto.EventFilter{}

	if value := c.Query("severity"); value != "" {
		severity, err := strconv.Atoi(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "severity must be a number",
			})
			return
		}

		filter.Severity = &severity
	}

	if value := c.Query("type_ulid"); value != "" {
		filter.TypeULID = &value
	}

	if value := c.Query("source_ulid"); value != "" {
		filter.SourceULID = &value
	}

	if value := c.Query("status"); value != "" {
		filter.Status = &value
	}

	if value := c.Query("from"); value != "" {
		filter.From = &value
	}

	if value := c.Query("to"); value != "" {
		filter.To = &value
	}

	if value := c.Query("hostname"); value != "" {
		filter.Hostname = &value
	}

	if value := c.Query("source_ip"); value != "" {
		filter.SourceIP = &value
	}

	events, err := h.service.GetAll(c.Request.Context(), filter)
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
