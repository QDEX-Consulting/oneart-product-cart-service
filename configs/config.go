package configs

import (
	"log"
	"os"
)

type Config struct {
	DBDSN     string
	Port      string
	JWTSecret string
}

func NewConfig() *Config {
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		dbDSN = "postgres://postgres:oneart-secret@35.244.41.139:5432/postgres?sslmode=disable"
		log.Printf("DB_DSN not sets, defaulting to: %s", dbDSN)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 for Cloud Run compatibility
		log.Printf("PORT not set, defaulting to: %s", port)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecret"
		log.Println("JWT_SECRET not set, defaulting to 'supersecret'")
	}

	return &Config{
		DBDSN:     dbDSN,
		Port:      port,
		JWTSecret: jwtSecret,
	}
}
