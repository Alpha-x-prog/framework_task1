package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/defects/app/backend/internal/repo"
)

type UsersHandler struct {
	DB *pgxpool.Pool
}

// GET /api/users?role=engineer&limit=100&offset=0
func (h *UsersHandler) List(c *gin.Context) {
	role := strings.TrimSpace(c.Query("role")) // engineer | manager | viewer | lead | ''
	limit := parseInt(c.Query("limit"), 100)
	offset := parseInt(c.Query("offset"), 0)

	// (опционально) валидация ролей, чтобы не мусорить БД
	if role != "" {
		switch role {
		case "engineer", "manager", "viewer", "lead":
			// ok
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown role"})
			return
		}
	}

	users, err := repo.ListUsers(c.Request.Context(), h.DB, role, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
