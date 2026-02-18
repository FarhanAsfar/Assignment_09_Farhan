package models

import "time"

// Airbnb property listing
type Listing struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Price     float64   `json:"price" db:"price"`
	Location  string    `json:"location" db:"location"`
	Rating    float64   `json:"rating" db:"rating"`
	URL       string    `json:"url" db:"url"`
	Bedrooms  int       `json:"bedrooms" db:"bedrooms"`
	Bathrooms int       `json:"bathrooms" db:"bathrooms"`
	Guests    int       `json:"guests" db:"guests"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// structure before normalization
type RawListing struct {
	Title     string
	Price     string
	Location  string
	Rating    string
	URL       string
	Bedrooms  int
	Bathrooms int
	Guests    int
}
