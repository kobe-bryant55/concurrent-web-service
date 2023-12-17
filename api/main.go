package main

import (
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/database"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/shared/config"
)

func main() {
	// Initialize application configs.
	cfg := config.Init()

	// Create a Postgres store.
	store := database.NewPostgresStore(database.WithUser(cfg.DB.User), database.WithName(cfg.DB.Name), database.WithPassword(cfg.DB.Password))

	// Initialize the Postgres store.
	store.InitDB()

	// Start server
	server := NewServer(":3333", cfg, store.GetInstance())
	server.Run()
}
