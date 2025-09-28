package config

import (
	"log"
	"os"
	"strings"

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

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("env %s is required", k)
	}
	return v
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
