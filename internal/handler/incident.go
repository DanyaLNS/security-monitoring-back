package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"security-monitor/internal/dto"
	"security-monitor/internal/service"
)

type IncidentHandler struct {
	service *service.IncidentService
}

func NewIncidentHandler(
	service *service.IncidentService,
) *IncidentHandler {
	return &IncidentHandler{service: service}
}

func (h *IncidentHandler) Create(c *gin.Context) {
	var input dto.CreateIncidentRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	incident, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

func (h *IncidentHandler) GetAll(c *gin.Context) {
	incidents, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

func (h *IncidentHandler) GetByULID(c *gin.Context) {
	ulid := c.Param("ulid")

	incident, err := h.service.GetByULID(c.Request.Context(), ulid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func (h *IncidentHandler) Update(c *gin.Context) {
	ulid := c.Param("ulid")

	var input dto.UpdateIncidentRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	incident, err := h.service.Update(c.Request.Context(), ulid, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func (h *IncidentHandler) Delete(c *gin.Context) {
	ulid := c.Param("ulid")

	if err := h.service.Delete(c.Request.Context(), ulid); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *IncidentHandler) GetIncidentEvents(c *gin.Context) {
	ulid := c.Param("ulid")

	events, err := h.service.GetEvents(c.Request.Context(), ulid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *IncidentHandler) AddEvents(c *gin.Context) {
	ulid := c.Param("ulid")

	var input dto.AddIncidentEventsRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	incident, err := h.service.AddEvents(c.Request.Context(), ulid, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, incident)
}
