package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kasteli-dev/tvtrackr-go/internal/adapters"
	"github.com/kasteli-dev/tvtrackr-go/internal/config"
	"github.com/kasteli-dev/tvtrackr-go/internal/core"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize PostgreSQL connection
	db, err := adapters.NewPostgresDB()
	if err != nil {
		log.Fatalf("Error connecting to Postgres: %v", err)
	}
	defer db.Close()

	service := core.NewUserServicePostgres(db.Pool)

	apiKey := os.Getenv("THETVDB_APIKEY")
	if apiKey == "" {
		log.Fatal("THETVDB_APIKEY environment variable not set")
	}
	tvdbClient := adapters.NewTheTVDBClient(apiKey)

	httpServer := adapters.NewHTTPServer(service, tvdbClient)
	router := httpServer.RegisterRoutes()

	log.Printf("Starting server on %s...", config.DefaultPort)
	err = http.ListenAndServe(config.DefaultPort, router)
	if err != nil {
		log.Fatal(err)
	}
}
