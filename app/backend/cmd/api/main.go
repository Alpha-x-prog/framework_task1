package main

import (
	"context"
	"example/defects/app/backend/internal/config"
	"log"

	"example/defects/app/backend/internal/db"
	httpx "example/defects/app/backend/internal/http"
	"example/defects/app/backend/internal/migrate"
)

func main() {
	// 1) читаем .env / переменные окружения
	cfg := config.Load()

	// 1) применяем миграции (up-only)
	if err := migrate.Up(cfg.DBURL); err != nil {
		log.Fatal("migrations failed:", err)
	}

	// 2) подключаемся к Postgres (пул соединений)
	pool, err := db.NewPool(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatal("db connection failed:", err)
	}
	defer pool.Close()

	// 3) собираем роутер (CORS, маршруты) и запускаем сервер
	r := httpx.NewRouter(httpx.Deps{
		DB:          pool,
		JWTSecret:   cfg.JWT,
		CORSOrigins: cfg.CORS,
	})
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatal("server error:", err)
	}
}
