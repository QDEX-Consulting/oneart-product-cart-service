package repository

import (
	"database/sql"
	"fmt"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/domain"
)

type CartRepository interface {
	AddToCart(item domain.CartItem) error
	GetCartItems(userID int64) ([]domain.CartItem, error)
	UpdateCartItem(item domain.CartItem) error
	RemoveCartItem(userID, productID int64) error
	ClearCart(userID int64) error
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddToCart(item domain.CartItem) error {
	query := `INSERT INTO cart_items (user_id, product_id, quantity, created_at, updated_at) 
              VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := r.db.Exec(query, item.UserID, item.ProductID, item.Quantity)
	return err
}

func (r *cartRepository) GetCartItems(userID int64) ([]domain.CartItem, error) {
	query := `SELECT id, user_id, product_id, quantity, created_at, updated_at 
              FROM cart_items 
              WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetCartItems query error: %v", err)
	}
	defer rows.Close()

	var items []domain.CartItem
	for rows.Next() {
		var ci domain.CartItem
		if err := rows.Scan(
			&ci.ID, &ci.UserID, &ci.ProductID, &ci.Quantity,
			&ci.CreatedAt, &ci.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, ci)
	}
	return items, nil
}

func (r *cartRepository) UpdateCartItem(item domain.CartItem) error {
	query := `UPDATE cart_items 
              SET quantity = $1, updated_at = NOW() 
              WHERE user_id = $2 AND product_id = $3`
	_, err := r.db.Exec(query, item.Quantity, item.UserID, item.ProductID)
	return err
}

func (r *cartRepository) RemoveCartItem(userID, productID int64) error {
	query := `DELETE FROM cart_items
              WHERE user_id = $1 AND product_id = $2`
	_, err := r.db.Exec(query, userID, productID)
	return err
}

func (r *cartRepository) ClearCart(userID int64) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
