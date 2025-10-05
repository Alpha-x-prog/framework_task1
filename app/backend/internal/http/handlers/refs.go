package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/defects/app/backend/internal/repo"
)

type RefsHandler struct{ DB *pgxpool.Pool }

func (h *RefsHandler) Statuses(c *gin.Context) {
	items, err := repo.ListStatuses(c, h.DB)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
func (h *RefsHandler) Roles(c *gin.Context) {
	items, err := repo.ListRoles(c, h.DB)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
