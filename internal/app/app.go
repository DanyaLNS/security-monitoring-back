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
	Incidents  *repository.IncidentRepo
}

type Services struct {
	Events     *service.EventService
	Sources    *service.SourceService
	EventTypes *service.EventTypeService
	Incidents  *service.IncidentService
}

type Handlers struct {
	Events     *handler.EventHandler
	Sources    *handler.SourceHandler
	EventTypes *handler.EventTypeHandler
	Incidents  *handler.IncidentHandler
}

func New(cfg *config.Config, db *pgxpool.Pool) *App {
	repositories := newRepositories(db)
	services := newServices(repositories)
	handlers := newHandlers(services)

	r := router.Init(
		handlers.Events,
		handlers.Sources,
		handlers.EventTypes,
		handlers.Incidents,
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
		Incidents:  repository.NewIncidentRepo(db),
	}
}

func newServices(repos *Repositories) *Services {
	return &Services{
		Events:     service.NewEventService(repos.Events),
		Sources:    service.NewSourceService(repos.Sources),
		EventTypes: service.NewEventTypeService(repos.EventTypes),
		Incidents:  service.NewIncidentService(repos.Incidents),
	}
}

func newHandlers(services *Services) *Handlers {
	return &Handlers{
		Events:     handler.NewEventHandler(services.Events),
		Sources:    handler.NewSourceHandler(services.Sources),
		EventTypes: handler.NewEventTypeHandler(services.EventTypes),
		Incidents:  handler.NewIncidentHandler(services.Incidents),
	}
}
