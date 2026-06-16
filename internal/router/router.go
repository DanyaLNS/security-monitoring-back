package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"security-monitor/internal/handler"
)

func Init(
	eventHandler *handler.EventHandler,
	sourceHandler *handler.SourceHandler,
	eventTypeHandler *handler.EventTypeHandler,
	incidentHandler *handler.IncidentHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	registerHealthRoutes(r)
	registerIngestRoutes(r, eventHandler)
	registerCoreRoutes(
		r,
		eventHandler,
		sourceHandler,
		eventTypeHandler,
		incidentHandler,
	)

	return r
}

func registerHealthRoutes(r *gin.Engine) {
	r.GET("/health", handler.HealthCheck)
}

func registerIngestRoutes(
	r *gin.Engine,
	eventHandler *handler.EventHandler,
) {
	ingest := r.Group("/api/ingest/v1")
	{
		ingest.GET("/health", handler.HealthCheck)

		events := ingest.Group("/events")
		{
			events.POST("", eventHandler.Create)
		}
	}
}

func registerCoreRoutes(
	r *gin.Engine,
	eventHandler *handler.EventHandler,
	sourceHandler *handler.SourceHandler,
	eventTypeHandler *handler.EventTypeHandler,
	incidentHandler *handler.IncidentHandler,
) {
	core := r.Group("/api/core/v1")
	{
		core.GET("/health", handler.HealthCheck)

		events := core.Group("/events")
		{
			events.GET("", eventHandler.GetAll)
			events.DELETE("/:ulid", eventHandler.Delete)
		}

		sources := core.Group("/sources")
		{
			sources.POST("", sourceHandler.Create)
			sources.GET("", sourceHandler.GetAll)
			sources.GET("/:ulid", sourceHandler.GetByULID)
			sources.DELETE("/:ulid", sourceHandler.Delete)
		}

		eventTypes := core.Group("/event-types")
		{
			eventTypes.POST("", eventTypeHandler.Create)
			eventTypes.GET("", eventTypeHandler.GetAll)
			eventTypes.GET("/:ulid", eventTypeHandler.GetByULID)
			eventTypes.DELETE("/:ulid", eventTypeHandler.Delete)
		}

		incidents := core.Group("/incidents")
		{
			incidents.POST("", incidentHandler.Create)
			incidents.GET("", incidentHandler.GetAll)
			incidents.GET("/:ulid", incidentHandler.GetByULID)
			incidents.PATCH("/:ulid", incidentHandler.Update)
			incidents.DELETE("/:ulid", incidentHandler.Delete)
			incidents.GET("/:ulid/events", incidentHandler.GetIncidentEvents)
			incidents.POST("/:ulid/events", incidentHandler.AddEvents)
		}
	}
}
