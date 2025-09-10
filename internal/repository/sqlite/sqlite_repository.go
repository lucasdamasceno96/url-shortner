package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lucasdamasceno96/url-shortner/internal/domain"
	"github.com/lucasdamasceno96/url-shortner/internal/repository"

	// The blank import is used for the side-effect of registering the driver.
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteRepository is the SQLite implementation of the URLRepository interface.
// It holds a reference to the database connection pool.
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new instance of SQLiteRepository.
// It opens a database connection and ensures the necessary table exists.
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	log.Println("Initializing SQLite repository...")

	// sql.Open() doesn't immediately create a connection. It prepares a connection pool.
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Ping the database to verify the connection is alive.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite database: %w", err)
	}

	// Create the 'urls' table if it doesn't already exist.
	if err := createTable(db); err != nil {
		return nil, err // The error is already descriptive from createTable
	}

	log.Println("SQLite repository initialized successfully.")
	return &SQLiteRepository{db: db}, nil
}

func createTable(db *sql.DB) error {
	const query = `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_url TEXT NOT NULL,
		short_code TEXT NOT NULL UNIQUE,
		created_at DATETIME NOT NULL
	);`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create 'urls' table: %w", err)
	}
	return nil
}

// Save persists a new ShortURL to the database.
func (r *SQLiteRepository) Save(ctx context.Context, url *domain.ShortURL) error {
	log.Printf("Attempting to save URL: %s with code: %s", url.OriginalURL, url.ShortCode)

	const query = "INSERT INTO urls (original_url, short_code, created_at) VALUES (?, ?, ?)"
	url.CreatedAt = time.Now()

	res, err := r.db.ExecContext(ctx, query, url.OriginalURL, url.ShortCode, url.CreatedAt)
	if err != nil {
		log.Printf("ERROR: Failed to save URL: %v", err)
		return fmt.Errorf("database error on save: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("ERROR: Failed to get last insert ID: %v", err)
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	url.ID = id
	log.Printf("URL saved successfully with ID: %d", url.ID)
	return nil
}

// FindByCode retrieves a ShortURL from the database by its short code.
func (r *SQLiteRepository) FindByCode(ctx context.Context, code string) (*domain.ShortURL, error) {
	log.Printf("Attempting to find URL by code: %s", code)

	const query = "SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = ?"

	var u domain.ShortURL
	row := r.db.QueryRowContext(ctx, query, code)

	err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("URL not found for code: %s", code)
			return nil, repository.ErrURLNotFound
		}
		log.Printf("ERROR: Failed to find URL by code: %v", err)
		return nil, fmt.Errorf("database error on find: %w", err)
	}

	log.Printf("URL found for code %s: (ID: %d)", code, u.ID)
	return &u, nil
}
