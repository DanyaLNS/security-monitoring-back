package router

import (
	"github.com/gin-gonic/gin"

	"security-monitor/internal/handler"
)

func Init(
	eventHandler *handler.EventHandler,
) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api/v1")

	events := api.Group("/events")
	{
		events.POST("", eventHandler.Create)
		events.GET("", eventHandler.GetAll)
		events.DELETE("/:ulid", eventHandler.Delete)
	}

	return r
}
