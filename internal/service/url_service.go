package service

import (
	"context"
	"fmt"
	"log"

	"github.com/lucasdamasceno96/url-shortner/internal/domain"
	"github.com/lucasdamasceno96/url-shortner/internal/repository"
	"github.com/lucasdamasceno96/url-shortner/internal/util"
)

// URLService contains the core business logic for URL shortening.
// It depends on a repository to persist data, abstracting away the database details.
type URLService struct {
	repo repository.URLRepository
}

// NewURLService creates a new instance of URLService.
// This is an example of Dependency Injection: we provide the service with its dependencies.
func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

// CreateShortURL handles the logic of creating, storing, and returning a short URL.
func (s *URLService) CreateShortURL(ctx context.Context, originalURL string) (*domain.ShortURL, error) {
	log.Printf("Service: Attempting to create short URL for: %s", originalURL)

	// In a real-world app, you'd also validate the URL format here.
	// For now, we keep it simple.

	shortCode := util.GenerateShortCode()
	// In a high-traffic system, you would need to check if this code already exists
	// and regenerate if it does. We'll omit that for simplicity.

	url := &domain.ShortURL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}

	err := s.repo.Save(ctx, url)
	if err != nil {
		log.Printf("ERROR: Service failed to save URL: %v", err)
		return nil, fmt.Errorf("could not save URL: %w", err)
	}

	log.Printf("Service: Successfully created short URL with code %s", url.ShortCode)
	return url, nil
}

// GetOriginalURL retrieves the original URL for a given short code.
func (s *URLService) GetOriginalURL(ctx context.Context, code string) (*domain.ShortURL, error) {
	log.Printf("Service: Attempting to find original URL for code: %s", code)

	url, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		log.Printf("ERROR: Service failed to find URL by code %s: %v", code, err)
		// We return the original error from the repository to allow the handler to act on it
		// (e.g., return a 404 if it's ErrURLNotFound).
		return nil, err
	}

	log.Printf("Service: Found original URL for code %s", code)
	return url, nil
}
