package router

import (
	"security-monitor/internal/handler"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	{
		api.GET("/health", handler.HealthCheck)
	}

	return r
}
