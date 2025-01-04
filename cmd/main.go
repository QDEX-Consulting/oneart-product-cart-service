package main

import (
	"context"
	"log"
	"net/http"

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

	// Connect to DB
	log.Println("Initializing database connection...")
	database, err := db.NewDB(cfg.DBDSN)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer database.Close()
	log.Println("Database connection established successfully!")

	// Initialize GCS client
	log.Println("Initializing GCS client...")
	ctx := context.Background()
	gcsClient, err := storage.NewGCSClient(ctx)
	if err != nil {
		log.Fatal("Error creating GCS client:", err)
	}
	log.Println("GCS client initialized successfully!")

	// Initialize repositories and services
	productRepo := repository.NewProductRepository(database)
	cartRepo := repository.NewCartRepository(database)
	productService := service.NewProductService(productRepo, gcsClient)
	cartService := service.NewCartService(cartRepo, productRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)

	// Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Setup router
	log.Println("Setting up router...")
	r := mux.NewRouter()

	// Product endpoints
	r.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	r.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	r.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/images/upload-temp", productHandler.UploadTempImage).Methods("POST")

	// Secure routes with authMiddleware
	r.Use(authMiddleware.JWTAuth)

	// Cart endpoints
	cartRouter := r.PathPrefix("/cart").Subrouter()
	cartRouter.Use(authMiddleware.JWTAuth)
	cartRouter.HandleFunc("", cartHandler.AddToCart).Methods("POST")
	cartRouter.HandleFunc("", cartHandler.GetCart).Methods("GET")

	// Start server
	log.Printf("Starting server on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
