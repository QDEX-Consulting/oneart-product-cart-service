package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	JWTSecret string
}

// NewAuthMiddleware creates a middleware with a known JWT secret.
func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{JWTSecret: secret}
}

// JWTAuth adapts to mux.MiddlewareFunc
func (am *AuthMiddleware) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// In a real system, also check signing method, e.g., HS256
			return []byte(am.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Extract user ID from claims (assuming "sub" is used for user ID)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userIDVal, ok := claims["sub"]
		if !ok {
			http.Error(w, "token missing 'sub' claim", http.StatusUnauthorized)
			return
		}

		// Typically "sub" might be a string; parse/convert as needed
		userID, ok := userIDVal.(float64)
		if !ok {
			http.Error(w, "invalid 'sub' claim type", http.StatusUnauthorized)
			return
		}

		// Store user ID in request context
		ctx := context.WithValue(r.Context(), "userID", int64(userID))

		// Pass to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
