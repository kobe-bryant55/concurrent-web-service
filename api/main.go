package main

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/database"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/shared/config"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/worker"
)

func main() {
	// Initialize application configs.
	cfg := config.Init()

	// Create a Postgres store.
	store := database.NewPostgresStore(database.WithUser(cfg.DB.User), database.WithName(cfg.DB.Name), database.WithPassword(cfg.DB.Password))

	// Initialize the Postgres store.
	store.InitDB()

	// Create a jobQueue.
	jobQueue := make(chan worker.Job, worker.MaxQueue)

	// Create and Start server.
	server := NewServer(":3333", cfg, store.GetInstance(), jobQueue)
	server.Run()
}
