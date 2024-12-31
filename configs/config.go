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
		dbDSN = "postgres://postgres:oneart-secret@34.124.225.173:5432/oneart_db?sslmode=disable"
		log.Printf("DB_DSN not set, defaulting to: %s", dbDSN)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // use 8082 if 8080 is used by identity-service
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
