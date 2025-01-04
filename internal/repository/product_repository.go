package repository

import (
	"database/sql"
	"fmt"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/domain"
)

type ProductRepository interface {
	ListProducts(filters map[string]interface{}) ([]domain.Product, error)
	GetProductByID(id int64) (*domain.Product, error)
	CreateProduct(product *domain.Product) error
	UpdateProduct(product *domain.Product) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// ListProducts retrieves products based on optional filters
func (r *productRepository) ListProducts(filters map[string]interface{}) ([]domain.Product, error) {
	query := `
        SELECT id, name, description, price, image_url, created_at, updated_at
        FROM products
        WHERE 1=1`

	var args []interface{}
	if category, ok := filters["category"]; ok {
		query += " AND category = $1"
		args = append(args, category)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}

// GetProductByID retrieves a single product by its ID
func (r *productRepository) GetProductByID(id int64) (*domain.Product, error) {
	query := `
        SELECT id, name, description, price, image_url, created_at, updated_at
        FROM products
        WHERE id = $1`

	var product domain.Product
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.CreatedAt, &product.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // No product found
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving product: %w", err)
	}
	return &product, nil
}

// CreateProduct inserts a new product into the database
func (r *productRepository) CreateProduct(product *domain.Product) error {
	query := `
        INSERT INTO products (name, description, price, image_url, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, product.Name, product.Description, product.Price, product.ImageURL).
		Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

// UpdateProduct updates an existing product in the database
func (r *productRepository) UpdateProduct(product *domain.Product) error {
	query := `
        UPDATE products
        SET name = $1, description = $2, price = $3, image_url = $4, updated_at = NOW()
        WHERE id = $5`
	_, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.ImageURL, product.ID)
	return err
}
