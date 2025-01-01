package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/QDEX-Core/oneart-product-cart-service/configs"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/db"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/handler"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/middleware"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/repository"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/service"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/storage"
)

func main() {
	// Load config
	cfg := configs.NewConfig()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT not set, defaulting to: %s", port)
	}

	// Connect to DB
	database, err := db.NewDB(cfg.DBDSN)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer database.Close()

	// Initialize repositories
	productRepo := repository.NewProductRepository(database)
	cartRepo := repository.NewCartRepository(database)

	// Initialize GCS client
	ctx := context.Background()
	gcsClient, err := storage.NewGCSClient(ctx)
	if err != nil {
		log.Fatal("Error creating GCS client:", err)
	}

	// Pass productRepo & gcsClient to productService
	productService := service.NewProductService(productRepo, gcsClient)
	cartService := service.NewCartService(cartRepo, productRepo) // unchanged

	// Handlers
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)

	// Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret) // or "supersecret"

	// Router
	r := mux.NewRouter()

	// Public product endpoints
	r.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	r.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")

	// NEW: Upload image endpoint (still "public" or you can add auth if you want)
	r.HandleFunc("/products/{id}/upload-image", productHandler.UploadImage).Methods("POST")

	// Protected cart endpoints
	cartRouter := r.PathPrefix("/cart").Subrouter()
	cartRouter.Use(authMiddleware.JWTAuth)
	cartRouter.HandleFunc("", cartHandler.AddToCart).Methods("POST")
	cartRouter.HandleFunc("", cartHandler.GetCart).Methods("GET")
	cartRouter.HandleFunc("/{productID}", cartHandler.UpdateCartItem).Methods("PUT")
	cartRouter.HandleFunc("/{productID}", cartHandler.RemoveCartItem).Methods("DELETE")
	cartRouter.HandleFunc("", cartHandler.ClearCart).Methods("DELETE")

	// Start server
	log.Printf("Product-Cart Service running on :%s\n", cfg.Port)
	if err = http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
