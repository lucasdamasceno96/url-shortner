package repository

import (
	"context"
	"errors"

	"github.com/lucasdamasceno96/url-shortner/internal/domain"
)

// ErrURLNotFound is a custom error returned when a short URL is not found.
// This allows the service layer to check specifically for this error type.
var ErrURLNotFound = errors.New("short URL not found")

// URLRepository is an interface that defines the contract for URL persistence.
// This decouples the business logic (service) from the data access layer (repository implementation).
type URLRepository interface {
	// Save persists a new ShortURL to the data store.
	// It takes a pointer to a ShortURL struct, which might be updated (e.g., with an ID) by the implementation.
	Save(ctx context.Context, url *domain.ShortURL) error

	// FindByCode retrieves a ShortURL from the data store by its short code.
	// It returns a pointer to the found ShortURL or ErrURLNotFound if it doesn't exist.
	FindByCode(ctx context.Context, code string) (*domain.ShortURL, error)
}
