package repository

import (
	"database/sql"
	"fmt"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/domain"
)

type ProductRepository interface {
	// Existing
	ListProducts(filter map[string]interface{}) ([]domain.Product, error)
	GetProductByID(id int64) (*domain.Product, error)

	// NEW: Must be declared in the interface
	UpdateProductImageURL(productID int64, imageURL string) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// ------------------------------------
// Implementation of interface methods
// ------------------------------------

// 1) ListProducts
func (r *productRepository) ListProducts(filter map[string]interface{}) ([]domain.Product, error) {
	query := `SELECT id, name, price, offer_price, category, country_of_origin, dimensions,
                     artist_name, image_url, created_at, updated_at
              FROM products
              WHERE 1=1`
	// We'll build filters dynamically (example: category)
	var args []interface{}

	if category, ok := filter["category"]; ok && category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	// Convert ? → $1, $2 for Postgres
	query = queryWithDollarParams(query)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ListProducts query error: %v", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.OfferPrice,
			&p.Category,
			&p.CountryOfOrigin,
			&p.Dimensions,
			&p.ArtistName,
			&p.ImageURL, // <--- scanning image_url
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// 2) GetProductByID
func (r *productRepository) GetProductByID(id int64) (*domain.Product, error) {
	query := `SELECT id, name, price, offer_price, category, country_of_origin,
                     dimensions, artist_name, image_url, created_at, updated_at
              FROM products
              WHERE id = $1`

	var p domain.Product
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.OfferPrice,
		&p.Category,
		&p.CountryOfOrigin,
		&p.Dimensions,
		&p.ArtistName,
		&p.ImageURL, // <--- scanning image_url
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// 3) UpdateProductImageURL (NEW)
func (r *productRepository) UpdateProductImageURL(productID int64, imageURL string) error {
	query := `
        UPDATE products
        SET image_url = $1, updated_at = NOW()
        WHERE id = $2
    `
	_, err := r.db.Exec(query, imageURL, productID)
	return err
}

// Helper to replace '?' with '$1' etc. in queries
func queryWithDollarParams(query string) string {
	var count int
	out := []rune{}
	for _, r := range query {
		if r == '?' {
			count++
			placeholder := fmt.Sprintf("$%d", count)
			for _, rr := range placeholder {
				out = append(out, rr)
			}
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}
