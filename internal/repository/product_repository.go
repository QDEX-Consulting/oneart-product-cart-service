package repository

import (
	"database/sql"
	"fmt"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/domain"
)

type ProductRepository interface {
	// Existing methods
	ListProducts(filter map[string]interface{}) ([]domain.Product, error)
	GetProductByID(id int64) (*domain.Product, error)
	UpdateProductImageURL(productID int64, imageURL string) error

	// New method
	CreateProduct(product *domain.Product) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// ListProducts retrieves products based on filters.
func (r *productRepository) ListProducts(filter map[string]interface{}) ([]domain.Product, error) {
	query := `SELECT id, name, price, offer_price, category, country_of_origin, dimensions,
                     artist_name, image_url, created_at, updated_at
              FROM products
              WHERE 1=1`
	// Dynamically build filters
	var args []interface{}
	if category, ok := filter["category"]; ok && category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}
	query = queryWithDollarParams(query) // Convert '?' to '$1, $2...' for PostgreSQL

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
			&p.ImageURL,
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

// GetProductByID retrieves a product by its ID.
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
		&p.ImageURL,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateProductImageURL updates the image URL for a specific product.
func (r *productRepository) UpdateProductImageURL(productID int64, imageURL string) error {
	query := `
        UPDATE products
        SET image_url = $1, updated_at = NOW()
        WHERE id = $2
    `
	_, err := r.db.Exec(query, imageURL, productID)
	return err
}

// CreateProduct inserts a new product into the database.
func (r *productRepository) CreateProduct(product *domain.Product) error {
	query := `
        INSERT INTO products (name, description, price, quantity, category, image_url, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, product.Name, product.Description, product.Price, product.Quantity, product.Category, product.ImageURL).
		Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

// Helper to convert '?' placeholders into '$1', '$2' for PostgreSQL.
func queryWithDollarParams(query string) string {
	var count int
	out := []rune{}
	for _, r := range query {
		if r == '?' {
			count++
			placeholder := fmt.Sprintf("$%d", count)
			out = append(out, []rune(placeholder)...)
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}
