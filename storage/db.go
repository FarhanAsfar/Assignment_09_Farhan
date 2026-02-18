package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/farhanasfar/airbnb-market-scraping-system/models"
	_ "github.com/lib/pq"
)

// DB wraps the database connection
type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection and runs migrations
func NewDB(dsn string) (*DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Run migrations
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✓ Database connected and migrated successfully")
	return db, nil
}

// migrate runs the database migrations
func (db *DB) migrate() error {
	// Create tables
	if _, err := db.conn.Exec(CreateListingsTableSQL); err != nil {
		return fmt.Errorf("failed to create listings table: %w", err)
	}

	// Create triggers
	if _, err := db.conn.Exec(UpdateUpdatedAtTriggerSQL); err != nil {
		return fmt.Errorf("failed to create triggers: %w", err)
	}

	return nil
}

// InsertListing inserts a new listing or updates if URL already exists
func (db *DB) InsertListing(listing *models.Listing) error {
	query := `
		INSERT INTO listings (title, price, location, rating, url, bedrooms, bathrooms, guests)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (url) DO UPDATE SET
			title = EXCLUDED.title,
			price = EXCLUDED.price,
			location = EXCLUDED.location,
			rating = EXCLUDED.rating,
			bedrooms = EXCLUDED.bedrooms,
			bathrooms = EXCLUDED.bathrooms,
			guests = EXCLUDED.guests,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`

	err := db.conn.QueryRow(
		query,
		listing.Title,
		listing.Price,
		listing.Location,
		listing.Rating,
		listing.URL,
		listing.Bedrooms,
		listing.Bathrooms,
		listing.Guests,
	).Scan(&listing.ID)

	if err != nil {
		return fmt.Errorf("failed to insert listing: %w", err)
	}

	return nil
}

// GetAllListings retrieves all listings from the database
func (db *DB) GetAllListings() ([]models.Listing, error) {
	query := `
		SELECT id, title, price, location, rating, url, bedrooms, bathrooms, guests, created_at, updated_at
		FROM listings
		ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query listings: %w", err)
	}
	defer rows.Close()

	var listings []models.Listing
	for rows.Next() {
		var l models.Listing
		err := rows.Scan(
			&l.ID, &l.Title, &l.Price, &l.Location, &l.Rating,
			&l.URL, &l.Bedrooms, &l.Bathrooms, &l.Guests,
			&l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}
		listings = append(listings, l)
	}

	return listings, nil
}

// close the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
