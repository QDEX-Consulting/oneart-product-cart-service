package domain

import "time"

// Product defines the structure for a product
type Product struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Price           float64   `json:"price"`
	Description     string    `json:"description"`
	Quantity        int       `json:"quantity"`
	OfferPrice      float64   `json:"offer_price,omitempty"`
	Category        string    `json:"category"`
	CountryOfOrigin string    `json:"country_of_origin,omitempty"`
	Dimensions      string    `json:"dimensions,omitempty"`
	ArtistName      string    `json:"artist_name,omitempty"`
	ImageURL        string    `json:"image_url,omitempty"` // URL for the product's image
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
