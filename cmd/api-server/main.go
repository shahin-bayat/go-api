// cmd/api-server/main.go
package main

import (
	"fmt"
	"log"

	"github.com/shahin-bayat/go-api/cmd/api-server/config"
	"github.com/shahin-bayat/go-api/cmd/seed"
	"github.com/shahin-bayat/go-api/internal/api"
	"github.com/shahin-bayat/go-api/internal/store"
)

func main() {
	// Parse command line flags
	config := config.ParseFlags()

	// Initialize database
	store, err := store.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// Seed the database if necessary
	if config.Seed {
		fmt.Println("Seeding the db...")
		seed.SeedAccounts(store)
	}

	// Start the API server
	server := api.NewAPIServer(":8000", store)
	server.Run()
}
