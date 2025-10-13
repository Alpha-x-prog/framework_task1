package httpx

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/defects/app/backend/internal/http/handlers"
	mw "example/defects/app/backend/internal/http/mv" // <-- фикс: было mv
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

	// Health
	r.GET("/api/healthz", handlers.Health)

	// Auth (публично)
	authH := &handlers.AuthHandler{DB: d.DB, Secret: d.JWTSecret}
	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	// ПУБЛИЧНЫЕ справочники (нужны на экране регистрации)
	refs := &handlers.RefsHandler{DB: d.DB}
	r.GET("/api/refs/roles", refs.Roles)

	// ====== ЗАЩИЩЁННАЯ ЗОНА ======
	api := r.Group("/api", mw.AuthRequired(d.JWTSecret))
	{
		// Профиль
		api.GET("/me", func(c *gin.Context) {
			uid, _ := c.Get("uid")
			role, _ := c.Get("role")
			c.JSON(200, gin.H{"uid": uid, "role": role})
		})

		// Справочники (просмотр статусов — всем аутентифицированным)
		api.GET("/refs/statuses", refs.Statuses)

		// ----- ПРОЕКТЫ -----
		prj := &handlers.ProjectsHandler{DB: d.DB}
		// Просмотр проектов: всем ролям (manager, engineer, lead, viewer)
		api.GET("/projects", prj.List)
		// Создание проектов: только менеджер/руководитель
		api.POST("/projects", mw.RequireRoles("manager", "lead"), prj.Create)
		// Если хочешь разрешить менеджеру назначать задачи отдельным эндпоинтом — добавишь здесь.

		// ----- ДЕФЕКТЫ -----
		def := &handlers.DefectsHandler{DB: d.DB}
		// Просмотр дефектов: всем ролям
		api.GET("/defects", def.List)
		// Создание дефектов: инженер/менеджер
		api.POST("/defects", mw.RequireRoles("engineer", "manager"), def.Create)
		// Изменение статуса: инженер/менеджер
		api.PATCH("/defects/:id/status", mw.RequireRoles("engineer", "manager"), def.UpdateStatus)

		// ----- КОММЕНТАРИИ -----
		cmt := &handlers.CommentsHandler{DB: d.DB}
		// Читать комментарии: всем ролям
		api.GET("/defects/:id/comments", cmt.List)
		// Добавлять комментарии: инженер/менеджер
		api.POST("/defects/:id/comments", mw.RequireRoles("engineer", "manager"), cmt.Create)

		// ----- ВЛОЖЕНИЯ -----
		att := &handlers.AttachmentsHandler{DB: d.DB}
		// Загрузка файлов: инженер/менеджер
		api.POST("/defects/:id/attachments", mw.RequireRoles("engineer", "manager"), att.Upload)

		// ----- ОТЧЁТЫ -----
		rep := &handlers.ReportsHandler{DB: d.DB}
		// Формирование/скачивание отчёта: менеджер/руководитель
		// (просмотр отчётности для руководителей и “заказчиков” → роль viewer тоже допускаем)
		api.GET("/reports/summary.csv", mw.RequireRoles("manager", "lead", "viewer"), rep.SummaryCSV)
	}

	return r
}
