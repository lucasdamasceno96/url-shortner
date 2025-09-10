package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lucasdamasceno96/url-shortner/internal/repository"
	"github.com/lucasdamasceno96/url-shortner/internal/service"
)

// We define request and response structs to be explicit about our API contract.

// ShortenRequest defines the structure for the request body of the POST /shorten endpoint.
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse defines the structure for the response body of the POST /shorten endpoint.
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// URLHandler holds the dependencies for the URL HTTP handlers, primarily the service.
type URLHandler struct {
	service *service.URLService
}

// NewURLHandler creates a new instance of URLHandler.
func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

// RegisterRoutes sets up the routing for the URL shortener endpoints.
func (h *URLHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/shorten", h.ShortenURL).Methods("POST")
	router.HandleFunc("/{shortCode}", h.RedirectURL).Methods("GET")
}

// ShortenURL is the handler for creating a new short URL.
func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: Received request to shorten URL")
	var req ShortenRequest

	// Decode the incoming JSON payload into our struct.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		log.Println("ERROR: URL is empty in request")
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	// Call the service to create the short URL.
	shortURL, err := h.service.CreateShortURL(r.Context(), req.URL)
	if err != nil {
		log.Printf("ERROR: Service failed to create short URL: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Construct the full short URL for the response.
	// In a real app, the base URL should come from a config file.
	baseURL := "http://localhost:8080"
	resp := ShortenResponse{
		ShortURL: fmt.Sprintf("%s/%s", baseURL, shortURL.ShortCode),
	}

	// Set headers and encode the JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("ERROR: Failed to encode response: %v", err)
	}
	log.Println("Handler: Successfully created and returned short URL")
}

// RedirectURL is the handler for redirecting a short code to its original URL.
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	// Extract the short code from the URL path variables.
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]
	log.Printf("Handler: Received request to redirect for code: %s", shortCode)

	// Call the service to find the original URL.
	url, err := h.service.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			log.Printf("Handler: URL not found for code: %s", shortCode)
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		// For any other error, return a generic server error.
		log.Printf("ERROR: Service failed to find URL: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Perform the redirect.
	log.Printf("Handler: Redirecting code %s to %s", shortCode, url.OriginalURL)
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
