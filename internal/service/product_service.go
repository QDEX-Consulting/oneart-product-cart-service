package service

import (
	"context"
	"fmt"
	"strings"
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
	// Fetch the products from the repository
	products, err := s.repo.ListProducts(filters)
	if err != nil {
		return nil, err
	}

	// Generate signed URLs for each product's ImageURL
	for i, product := range products {
		if product.ImageURL != "" {
			// Extract object name from the URL if needed
			objectName := product.ImageURL
			if strings.HasPrefix(objectName, "https://storage.googleapis.com/") {
				objectName = strings.TrimPrefix(objectName, "https://storage.googleapis.com/"+s.gcsClient.GetBucketName()+"/")
			}

			// Generate a signed URL for the object
			signedURL, err := s.gcsClient.GenerateSignedURL(objectName, time.Hour)
			if err != nil {
				return nil, fmt.Errorf("failed to generate signed URL for product %d: %w", product.ID, err)
			}

			// Update the product's ImageURL with the signed URL
			products[i].ImageURL = signedURL
		}
	}

	return products, nil
}

func (s *productService) GetProduct(id int64) (*domain.Product, error) {
	// Fetch the product from the repository
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	// Ensure the ImageURL contains only the object name
	if product.ImageURL != "" {
		// Extract object name from the URL if needed
		objectName := product.ImageURL
		if strings.HasPrefix(objectName, "https://storage.googleapis.com/") {
			// Remove the base URL prefix
			objectName = strings.TrimPrefix(objectName, "https://storage.googleapis.com/"+s.gcsClient.GetBucketName()+"/")
		}

		// Generate a signed URL for the object
		signedURL, err := s.gcsClient.GenerateSignedURL(objectName, time.Hour)
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
