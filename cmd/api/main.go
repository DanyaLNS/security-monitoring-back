package main

import (
	"log"
	"security-monitor/internal/config"
	"security-monitor/internal/router"
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

	r := router.Init()

	log.Printf("server started on :%s", cfg.AppPort)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
