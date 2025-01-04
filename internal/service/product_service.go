package service

import (
	"context"
	"fmt"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/domain"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/repository"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/storage"
)

// ProductService interface defines the service layer contract
type ProductService interface {
	ListProducts(filters map[string]interface{}) ([]domain.Product, error)
	GetProduct(id int64) (*domain.Product, error)
	CreateProduct(product *domain.Product) error
	UploadProductImage(ctx context.Context, productID int64, fileName string, fileData []byte) (string, error)
}

// productService is the implementation of ProductService
type productService struct {
	repo      repository.ProductRepository
	gcsClient *storage.GCSClient
}

// NewProductService creates a new productService instance
func NewProductService(repo repository.ProductRepository, gcsClient *storage.GCSClient) ProductService {
	return &productService{
		repo:      repo,
		gcsClient: gcsClient,
	}
}

// ListProducts retrieves a list of products based on filters
func (s *productService) ListProducts(filters map[string]interface{}) ([]domain.Product, error) {
	return s.repo.ListProducts(filters)
}

// GetProduct retrieves a single product by ID
func (s *productService) GetProduct(id int64) (*domain.Product, error) {
	return s.repo.GetProductByID(id)
}

// CreateProduct validates and creates a new product
func (s *productService) CreateProduct(product *domain.Product) error {
	// Validate product fields
	if err := validateProduct(product); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Save the product to the database
	return s.repo.CreateProduct(product)
}

// UploadProductImage uploads an image for a product and updates its record in the database
func (s *productService) UploadProductImage(
	ctx context.Context,
	productID int64,
	fileName string,
	fileData []byte,
) (string, error) {
	// Validate inputs
	if productID <= 0 {
		return "", fmt.Errorf("invalid product ID")
	}
	if fileName == "" {
		return "", fmt.Errorf("file name is required")
	}
	if len(fileName) > 255 {
		return "", fmt.Errorf("file name must not exceed 255 characters")
	}
	if len(fileData) == 0 {
		return "", fmt.Errorf("file data cannot be empty")
	}

	// Ensure product exists
	_, err := s.repo.GetProductByID(productID)
	if err != nil {
		return "", fmt.Errorf("product with ID %d does not exist", productID)
	}

	// Upload image to GCS and get the URL
	objectName := fmt.Sprintf("products/%d/%s", productID, fileName)
	url, err := s.gcsClient.UploadImage(ctx, objectName, fileData)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	// Update the product record with the image URL
	err = s.repo.UpdateProductImageURL(productID, url)
	if err != nil {
		return "", fmt.Errorf("failed to update image URL in database: %w", err)
	}

	return url, nil
}

// validateProduct performs validation on a product before saving it
func validateProduct(product *domain.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if len(product.Name) > 100 {
		return fmt.Errorf("product name must not exceed 100 characters")
	}
	if product.Price <= 0 {
		return fmt.Errorf("product price must be greater than zero")
	}
	if product.Quantity < 0 {
		return fmt.Errorf("product quantity cannot be negative")
	}
	if product.Category == "" {
		return fmt.Errorf("product category is required")
	}
	if len(product.Description) > 500 {
		return fmt.Errorf("product description must not exceed 500 characters")
	}
	if product.ImageURL != "" && len(product.ImageURL) > 2048 {
		return fmt.Errorf("image URL must not exceed 2048 characters")
	}

	return nil
}
