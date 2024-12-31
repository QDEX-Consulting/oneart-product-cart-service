package service

import (
	"errors"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/domain"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/repository"
)

type CartService interface {
	AddToCart(userID, productID int64, quantity int) error
	GetCart(userID int64) ([]domain.CartItem, error)
	UpdateCartItem(userID, productID int64, quantity int) error
	RemoveCartItem(userID, productID int64) error
	ClearCart(userID int64) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartService) AddToCart(userID, productID int64, quantity int) error {
	// Check if product exists
	product, err := s.productRepo.GetProductByID(productID)
	if err != nil || product == nil {
		return errors.New("product not found")
	}
	// Additional checks? e.g. stock, max quantity, etc.

	cartItem := domain.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.cartRepo.AddToCart(cartItem)
}

func (s *cartService) GetCart(userID int64) ([]domain.CartItem, error) {
	return s.cartRepo.GetCartItems(userID)
}

func (s *cartService) UpdateCartItem(userID, productID int64, quantity int) error {
	// Could also check if product exists, or if quantity is valid
	cartItem := domain.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.cartRepo.UpdateCartItem(cartItem)
}

func (s *cartService) RemoveCartItem(userID, productID int64) error {
	return s.cartRepo.RemoveCartItem(userID, productID)
}

func (s *cartService) ClearCart(userID int64) error {
	return s.cartRepo.ClearCart(userID)
}
