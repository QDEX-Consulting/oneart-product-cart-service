package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/QDEX-Core/oneart-product-cart-service/internal/service"
	"github.com/gorilla/mux"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cs service.CartService) *CartHandler {
	return &CartHandler{cartService: cs}
}

// GetCart handles GET /cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the JWT middleware context
	ctxUserID := r.Context().Value("userID")
	userID, ok := ctxUserID.(int64)
	if !ok {
		http.Error(w, "User ID not found in token context", http.StatusUnauthorized)
		return
	}

	items, err := h.cartService.GetCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

// UpdateCartItem handles PUT /cart/{productID}
func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr := vars["productID"]
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid productID", http.StatusBadRequest)
		return
	}

	// Extract userID from the JWT middleware context
	ctxUserID := r.Context().Value("userID")
	userID, ok := ctxUserID.(int64)
	if !ok {
		http.Error(w, "User ID not found in token context", http.StatusUnauthorized)
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.cartService.UpdateCartItem(userID, productID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RemoveCartItem handles DELETE /cart/{productID}
func (h *CartHandler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr := vars["productID"]
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid productID", http.StatusBadRequest)
		return
	}

	// Extract userID from the JWT middleware context
	ctxUserID := r.Context().Value("userID")
	userID, ok := ctxUserID.(int64)
	if !ok {
		http.Error(w, "User ID not found in token context", http.StatusUnauthorized)
		return
	}

	err = h.cartService.RemoveCartItem(userID, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ClearCart handles DELETE /cart
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the JWT middleware context
	ctxUserID := r.Context().Value("userID")
	userID, ok := ctxUserID.(int64)
	if !ok {
		http.Error(w, "User ID not found in token context", http.StatusUnauthorized)
		return
	}

	err := h.cartService.ClearCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AddToCart handles POST /cart
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the JWT middleware context
	ctxUserID := r.Context().Value("userID")
	userID, ok := ctxUserID.(int64)
	if !ok {
		http.Error(w, "User ID not found in token context", http.StatusUnauthorized)
		return
	}

	// The request body might only contain productID, quantity
	var req struct {
		ProductID int64 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.cartService.AddToCart(userID, req.ProductID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
