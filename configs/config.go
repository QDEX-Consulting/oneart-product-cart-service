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
		dbDSN = "postgres://postgres:oneart-secret@/postgres?host=/cloudsql/qdex-401002:asia-south1:oneart-postgres"
		log.Printf("DB_DSN not set, defaulting to: %s", dbDSN)
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
