package main

import (
	_ "github.com/MehmetTalhaSeker/concurrent-web-service/docs"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/database"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/shared/config"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/worker"
)

// @title API
// @version 0.01
// @description Concurrent Data Processing Web Service

// @host localhost:3333
// @BasePath /
func main() {
	// Initialize application configs.
	cfg := config.Init()

	// Create a Postgres store.
	store := database.NewPostgresStore(database.WithUser(cfg.DB.User), database.WithName(cfg.DB.Name), database.WithPassword(cfg.DB.Password))

	// Initialize the Postgres store.
	store.InitDB()

	// Create a jobQueue.
	jobQueue := make(chan worker.Job, 10)

	// Create and Start server.
	server := NewServer(":3333", cfg, store.GetInstance(), jobQueue)
	server.Run()
}
