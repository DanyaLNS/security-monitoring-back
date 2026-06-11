package app

import (
	"security-monitor/internal/config"
	"security-monitor/internal/handler"
	"security-monitor/internal/repository"
	"security-monitor/internal/router"
	"security-monitor/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg    *config.Config
	db     *pgxpool.Pool
	router *gin.Engine
}

type Repositories struct {
	Events     *repository.EventRepo
	Sources    *repository.EventSourceRepo
	EventTypes *repository.EventTypeRepo
	Analysis   *repository.EventAnalysisRepo
}

type Services struct {
	Events     *service.EventService
	Sources    *service.SourceService
	EventTypes *service.EventTypeService
	Dashboard  *service.DashboardService
	Incidents  *service.IncidentService
	Analysis   *service.AnalysisService
}

type Handlers struct {
	Events     *handler.EventHandler
	Sources    *handler.SourceHandler
	EventTypes *handler.EventTypeHandler
	Dashboard  *handler.DashboardHandler
	Incidents  *handler.IncidentHandler
	Analysis   *handler.AnalysisHandler
}

func New(cfg *config.Config, db *pgxpool.Pool) *App {
	repositories := newRepositories(db)
	services := newServices(repositories)
	handlers := newHandlers(services)

	r := router.Init(
		handlers.Events,
		handlers.Sources,
		handlers.EventTypes,
		handlers.Dashboard,
		handlers.Incidents,
		handlers.Analysis,
	)

	return &App{
		cfg:    cfg,
		db:     db,
		router: r,
	}
}

func (a *App) Run() error {
	return a.router.Run(":" + a.cfg.AppPort)
}

func newRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		Events:     repository.NewEventRepo(db),
		Sources:    repository.NewEventSourceRepo(db),
		EventTypes: repository.NewEventTypeRepo(db),
		Analysis:   repository.NewEventAnalysisRepo(db),
	}
}

func newServices(repos *Repositories) *Services {
	return &Services{
		Events: service.NewEventService(repos.Events),

		// Пока эти сервисы у тебя заглушечные.
		// Позже лучше заменить на конструкторы вида:
		// service.NewSourceService(repos.Sources)
		// service.NewEventTypeService(repos.EventTypes)
		// service.NewAnalysisService(repos.Analysis)
		Sources:    &service.SourceService{},
		EventTypes: &service.EventTypeService{},
		Dashboard:  &service.DashboardService{},
		Incidents:  &service.IncidentService{},
		Analysis:   &service.AnalysisService{},
	}
}

func newHandlers(services *Services) *Handlers {
	return &Handlers{
		Events:     handler.NewEventHandler(services.Events),
		Sources:    handler.NewSourceHandler(services.Sources),
		EventTypes: handler.NewEventTypeHandler(services.EventTypes),
		Dashboard:  handler.NewDashboardHandler(services.Dashboard),
		Incidents:  handler.NewIncidentHandler(services.Incidents),
		Analysis:   handler.NewAnalysisHandler(services.Analysis),
	}
}
