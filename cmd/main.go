package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/thedevflex/kubi8p/internal/cache"
	"github.com/thedevflex/kubi8p/internal/handler"
	"github.com/thedevflex/kubi8p/internal/k8utils"
	"github.com/thedevflex/kubi8p/internal/server"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v, Server will assume default values", err)
	}

	installerCache := cache.NewInstallerCache()
	admin, err := k8utils.NewAdmin()
	if err != nil {
		log.Fatalf("Failed to create admin: %v", err)
	}
	handler := handler.NewHandler(installerCache, admin)

	// Create a new server
	server := server.NewServer()
	server.Register(handler)
	server.Start("8080")
}
