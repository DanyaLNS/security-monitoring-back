package main

import (
	"log"
	"security-monitor/internal/config"
	"security-monitor/internal/handler"
	"security-monitor/internal/repository"
	"security-monitor/internal/router"
	"security-monitor/internal/service"
	"security-monitor/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewPostgresConnection(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewPostgresEventRepository(db)
	service := service.NewEventService(repo)
	handler := handler.NewEventHandler(service)

	r := router.Init(handler)

	log.Fatal(r.Run(":" + cfg.AppPort))
}
