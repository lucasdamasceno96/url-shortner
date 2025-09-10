package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lucasdamasceno96/url-shortner/internal/handler"
	"github.com/lucasdamasceno96/url-shortner/internal/repository/sqlite"
	"github.com/lucasdamasceno96/url-shortner/internal/service"
)

func main() {
	log.Println("Starting URL Shortener API...")

	// --- Dependency Injection ---
	// The layers are built from the inside out: repository -> service -> handler

	// 1. Create the data directory if it doesn't exist
	const dataDir = "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// 2. Initialize Repository
	// This is the persistence layer.
	dbPath := dataDir + "/shortener.db"
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// 3. Initialize Service
	// The service layer contains the business logic and uses the repository.
	urlService := service.NewURLService(repo)

	// 4. Initialize Handler
	// The handler layer is responsible for HTTP requests and uses the service.
	urlHandler := handler.NewURLHandler(urlService)

	// 5. Setup Router
	router := mux.NewRouter()
	urlHandler.RegisterRoutes(router)

	// 6. Start Server
	port := ":8080"
	log.Printf("Server is listening on port %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
