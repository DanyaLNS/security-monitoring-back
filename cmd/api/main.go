package main

import (
	"log"

	"security-monitor/internal/app"
	"security-monitor/internal/config"
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

	application := app.New(cfg, db)

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
