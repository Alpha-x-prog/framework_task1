package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr  string   // :8080
	DBURL string   // postgres://...
	JWT   string   // секрет для JWT
	CORS  []string // origin'ы для фронта
}

func Load() Config {
	_ = godotenv.Load(".env")

	addr := get("API_ADDR", ":8080")
	db := must("DB_URL")
	jwt := must("JWT_SECRET")
	cors := splitCSV(get("CORS_ORIGINS", "http://localhost:5173"))
	return Config{Addr: addr, DBURL: db, JWT: jwt, CORS: cors}
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
