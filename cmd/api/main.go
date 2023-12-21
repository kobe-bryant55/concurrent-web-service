package main

import (
	_ "github.com/MehmetTalhaSeker/concurrent-web-service/docs"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/database"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/logger"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/rba"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/shared/config"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/validatorutils"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/worker"
	"log"
)

// @title API
// @version 0.01
// @description Concurrent Data Processing Web Service

// @host localhost:3333
// @BasePath /
func main() {
	// New logger.
	logger.Init()

	// Initialize application configs.
	cfg, err := config.Init()
	if err != nil {
		logger.ErrorLog.Println(err)
		log.Fatal(err)
	}

	// Create a Postgres store.
	store, err := database.NewPostgresStore(database.WithUser(cfg.DB.User), database.WithName(cfg.DB.Name), database.WithPassword(cfg.DB.Password))
	if err != nil {
		logger.ErrorLog.Println(err)
		log.Fatal(err)
	}

	// Initialize the Postgres store.
	if err = store.InitDB(); err != nil {
		logger.ErrorLog.Println(err)
		log.Fatal(err)
	}

	// Create a jobQueue.
	jobQueue := make(chan worker.Job, cfg.Worker.MaxQueue)

	// Create RBA for role based authentication.
	rb := rba.New()

	vl := validatorutils.NewValidator()

	// Create and Start server.
	server := NewServer(":3333", cfg, store.GetInstance(), jobQueue, rb, vl)

	if err = server.Run(); err != nil {
		logger.ErrorLog.Println(err)
		log.Fatal(err)
	}
}
