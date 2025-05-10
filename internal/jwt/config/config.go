package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type JWT struct {
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
}

func LoadJWT() (*JWT, error) {
	err := godotenv.Load("./internal/orchestrator/.env")
	if err != nil {
		log.Fatalln("jwt/config/LoadJWT(): Не удалось загрузить .env файл", err)
	}
	jwt := &JWT{}
	jwt.Secret = Getenv("JWT_SECRET", "your-secure-secret-key")
	value, err := strconv.ParseInt(Getenv("TOKEN_EXPIRY", "30"), 10, 64)
	if err != nil {
		log.Fatalf("jwt/config/func LoadJWT(): что-то пошло не так: %q с парсингом TOKEN_EXPIRY: %d", err, value)
	}
	jwt.TokenExpiry = time.Duration(value) * time.Hour
	value, err = strconv.ParseInt(Getenv("REFRESH_EXPIRY", "30"), 10, 64)
	if err != nil {
		log.Fatalf("jwt/config/func LoadJWT(): что-то пошло не так: %q с парсингом REFRESH_EXPIRY: %d", err, value)
	}
	jwt.RefreshExpiry = time.Duration(value) * time.Hour
	return jwt, nil
}

func Getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
