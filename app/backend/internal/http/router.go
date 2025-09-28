package httpx

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"example/defects/app/backend/internal/http/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	DB          *pgxpool.Pool
	JWTSecret   string
	CORSOrigins []string
}

func NewRouter(d Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     d.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// публичный healthcheck
	r.GET("/api/healthz", handlers.Health)

	// здесь позже добавим auth/login и защищённые группы
	return r
}
