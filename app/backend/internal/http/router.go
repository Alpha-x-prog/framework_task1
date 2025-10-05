package httpx

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"example/defects/app/backend/internal/http/handlers"

	mw "example/defects/app/backend/internal/http/mv"

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

	// login
	authH := &handlers.AuthHandler{DB: d.DB, Secret: d.JWTSecret}

	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	// protected
	api := r.Group("/api", mw.AuthRequired(d.JWTSecret))
	{
		prj := &handlers.ProjectsHandler{DB: d.DB}
		api.GET("/projects", prj.List)
		api.POST("/projects", prj.Create)

		def := &handlers.DefectsHandler{DB: d.DB}
		api.GET("/defects", def.List)
		api.POST("/defects", def.Create)

		cmt := &handlers.CommentsHandler{DB: d.DB}
		api.GET("/defects/:id/comments", cmt.List)
		api.POST("/defects/:id/comments", cmt.Create)

		api.GET("/me", func(c *gin.Context) {
			uid, _ := c.Get("uid")
			role, _ := c.Get("role")
			c.JSON(200, gin.H{"uid": uid, "role": role})
		})
	}
	return r
}
