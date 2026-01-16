package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/intaek-h/ghofig/internal/model"
	_ "modernc.org/sqlite"
)

var db *sql.DB

// Init initializes the database from embedded bytes.
// It writes the embedded DB to a temp file and opens it.
func Init(embeddedDB []byte) error {
	// Write embedded DB to temp file
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "ghofig.db")

	if err := os.WriteFile(tempFile, embeddedDB, 0644); err != nil {
		return fmt.Errorf("failed to write temp db: %w", err)
	}

	var err error
	db, err = sql.Open("sqlite", tempFile+"?mode=ro")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close closes the database connection.
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// Search searches for configs matching the query.
// Results prioritize title matches over description matches.
func Search(query string) ([]model.Config, error) {
	if query == "" {
		return getAllConfigs()
	}

	likeQuery := "%" + query + "%"

	rows, err := db.Query(`
		SELECT id, title, description 
		FROM configs 
		WHERE title LIKE ? OR description LIKE ?
		ORDER BY 
			CASE WHEN title LIKE ? THEN 0 ELSE 1 END,
			title
		LIMIT 50
	`, likeQuery, likeQuery, likeQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanConfigs(rows)
}

// GetByID retrieves a single config by its ID.
func GetByID(id int) (*model.Config, error) {
	row := db.QueryRow("SELECT id, title, description FROM configs WHERE id = ?", id)

	var config model.Config
	err := row.Scan(&config.ID, &config.Title, &config.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("config not found: %d", id)
		}
		return nil, err
	}

	return &config, nil
}

// getAllConfigs returns all configs ordered by title.
func getAllConfigs() ([]model.Config, error) {
	rows, err := db.Query("SELECT id, title, description FROM configs ORDER BY title LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanConfigs(rows)
}

// scanConfigs scans rows into a slice of Config.
func scanConfigs(rows *sql.Rows) ([]model.Config, error) {
	var configs []model.Config
	for rows.Next() {
		var c model.Config
		if err := rows.Scan(&c.ID, &c.Title, &c.Description); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}
