package domain

import "time"

// ShortURL is the core struct for our application.
// It represents the mapping between an original URL and its shortened version.
type ShortURL struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
}
