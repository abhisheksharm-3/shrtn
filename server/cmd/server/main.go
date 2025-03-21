package main

import (
	"log"

	"github.com/abhisheksharm-3/shrtn/internal/api"
	"github.com/abhisheksharm-3/shrtn/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	r := api.SetupRouter(cfg)
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
