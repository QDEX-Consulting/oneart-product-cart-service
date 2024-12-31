package domain

import "time"

type Product struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	OfferPrice      float64 `json:"offer_price"`
	Category        string  `json:"category"`
	CountryOfOrigin string  `json:"country_of_origin"`
	Dimensions      string  `json:"dimensions"`
	ArtistName      string  `json:"artist_name"`
	// Add any other columns, e.g., seller_name if you have it
	ImageURL  string    `json:"image_url"` // <--- new field for image URL
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
