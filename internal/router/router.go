package router

import (
	"github.com/gin-gonic/gin"

	"security-monitor/internal/handler"
)

func Init(
	eventHandler *handler.EventHandler,
	sourceHandler *handler.SourceHandler,
	eventTypeHandler *handler.EventTypeHandler,
	dashboardHandler *handler.DashboardHandler,
	incidentHandler *handler.IncidentHandler,
	analysisHandler *handler.AnalysisHandler,
) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api/v1")

	events := api.Group("/events")
	{
		events.POST("", eventHandler.Create)
		events.GET("", eventHandler.GetAll)
		events.DELETE("/:ulid", eventHandler.Delete)
	}

	sources := api.Group("/sources")
	{
		sources.POST("", sourceHandler.Create)
		sources.GET("", sourceHandler.GetAll)
		sources.GET("/:ulid", sourceHandler.GetByULID)
		sources.DELETE("/:ulid", sourceHandler.Delete)
	}

	eventTypes := api.Group("/event-types")
	{
		eventTypes.POST("", eventTypeHandler.Create)
		eventTypes.GET("", eventTypeHandler.GetAll)
		eventTypes.GET("/:ulid", eventTypeHandler.GetByULID)
		eventTypes.DELETE("/:ulid", eventTypeHandler.Delete)
	}

	dashboard := api.Group("/dashboard")
	{
		dashboard.GET("/metrics", dashboardHandler.GetMetrics)
		dashboard.GET("/timeline", dashboardHandler.GetTimeline)
		dashboard.GET("/severity-distribution", dashboardHandler.GetSeverityDistribution)
		dashboard.GET("/top-sources", dashboardHandler.GetTopSources)
		dashboard.GET("/recent-events", dashboardHandler.GetRecentEvents)
	}

	incidents := api.Group("/incidents")
	{
		incidents.GET("", incidentHandler.GetAll)
		incidents.GET("/:ulid", incidentHandler.GetByULID)
		incidents.GET("/:ulid/events", incidentHandler.GetIncidentEvents)
	}

	analysis := api.Group("/analysis")
	{
		analysis.GET("", analysisHandler.GetAll)
		analysis.GET("/:event_ulid", analysisHandler.GetByEventULID)
	}

	return r
}
