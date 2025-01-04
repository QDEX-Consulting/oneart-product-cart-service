package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/domain"
	"github.com/QDEX-Core/oneart-product-cart-service/internal/service"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(ps service.ProductService) *ProductHandler {
	return &ProductHandler{productService: ps}
}

// ListProducts handles GET /products?category=painting&...
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	filters := map[string]interface{}{}
	category := r.URL.Query().Get("category")
	if category != "" {
		filters["category"] = category
	}

	products, err := h.productService.ListProducts(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(products)
}

// GetProduct handles GET /products/{id}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(product)
}

// NEW: UploadImage handles POST /products/{id}/upload-image
func (h *ProductHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr := vars["id"]
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// parse up to 10 MB of multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// "image" is the file field
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Missing 'image' file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// read file bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file bytes", http.StatusInternalServerError)
		return
	}

	// upload & update DB
	url, err := h.productService.UploadProductImage(r.Context(), productID, fileHeader.Filename, fileBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return JSON with the new image_url
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"image_url": url,
	})
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req domain.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.productService.CreateProduct(&req); err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(req)
}
