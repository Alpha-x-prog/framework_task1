package main

import (
	"context"
	"example/defects/app/backend/internal/config"

	"example/defects/app/backend/internal/db"
	httpx "example/defects/app/backend/internal/http"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1) грузим конфиг из .env / env
	cfg := config.Load()

	// 2) создаём пул соединений с БД
	pool, err := db.NewPool(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatal("db connection:", err)
	}
	defer pool.Close()

	// 3) собираем роутер (CORS, маршруты)
	r := httpx.NewRouter(httpx.Deps{
		DB:          pool,
		JWTSecret:   cfg.JWT,
		CORSOrigins: cfg.CORS,
	})

	// 4) http-сервер + graceful shutdown
	srv := &http.Server{Addr: cfg.Addr, Handler: r}

	go func() {
		log.Println("API listening on", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("server stopped")
}
