package service

import (
	"context"
	"fmt"
	"time"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/domain"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/repository"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/storage"
)

type ProductService interface {
	ListProducts(filters map[string]interface{}) ([]domain.Product, error)
	GetProduct(id int64) (*domain.Product, error)
	CreateProduct(product *domain.Product) error
	UpdateProduct(product *domain.Product) error
	UploadTempImage(ctx context.Context, fileName string, fileData []byte) (string, error)
}

type productService struct {
	repo      repository.ProductRepository
	gcsClient *storage.GCSClient
}

func NewProductService(repo repository.ProductRepository, gcsClient *storage.GCSClient) ProductService {
	return &productService{
		repo:      repo,
		gcsClient: gcsClient,
	}
}

func (s *productService) ListProducts(filters map[string]interface{}) ([]domain.Product, error) {
	return s.repo.ListProducts(filters)
}

func (s *productService) GetProduct(id int64) (*domain.Product, error) {
	// Fetch the product from the repository
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	// Generate a signed URL for the product image if an image URL exists
	if product.ImageURL != "" {
		signedURL, err := s.gcsClient.GenerateSignedURL(product.ImageURL, time.Hour)
		if err != nil {
			return nil, fmt.Errorf("failed to generate signed URL: %w", err)
		}
		product.ImageURL = signedURL
	}

	return product, nil
}

func (s *productService) CreateProduct(product *domain.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.CreateProduct(product)
}

func (s *productService) UpdateProduct(product *domain.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.UpdateProduct(product)
}

func (s *productService) UploadTempImage(ctx context.Context, fileName string, fileData []byte) (string, error) {
	objectName := fmt.Sprintf("temp/%s", fileName)
	return s.gcsClient.UploadImage(ctx, objectName, fileData)
}

func validateProduct(product *domain.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if product.Price <= 0 {
		return fmt.Errorf("product price must be greater than zero")
	}
	return nil
}
