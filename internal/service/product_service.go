package service

import (
	"context"
	"fmt"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/domain"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/repository"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/storage"
)

type ProductService interface {
	ListProducts(filters map[string]interface{}) ([]domain.Product, error)
	GetProduct(id int64) (*domain.Product, error)

	// NEW: Allows uploading an image for a product
	UploadProductImage(ctx context.Context, productID int64, fileName string, fileData []byte) (string, error)
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
	return s.repo.GetProductByID(id)
}

// NEW: UploadProductImage
func (s *productService) UploadProductImage(
	ctx context.Context,
	productID int64,
	fileName string,
	fileData []byte,
) (string, error) {

	// e.g. objectName = "products/123/image.jpg"
	objectName := fmt.Sprintf("products/%d/%s", productID, fileName)

	// 1. Upload to GCS → returns a GCS URL
	url, err := s.gcsClient.UploadImage(ctx, objectName, fileData)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to GCS: %w", err)
	}

	// 2. Update DB (image_url column)
	err = s.repo.UpdateProductImageURL(productID, url)
	if err != nil {
		return "", fmt.Errorf("failed to update image URL in DB: %w", err)
	}

	return url, nil
}
