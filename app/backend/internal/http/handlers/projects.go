package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/defects/app/backend/internal/repo"
)

type ProjectsHandler struct {
	DB *pgxpool.Pool
}

func (h *ProjectsHandler) List(c *gin.Context) {
	items, err := repo.ListProjects(c, h.DB)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ProjectsHandler) Create(c *gin.Context) {
	var req struct {
		Name     string  `json:"name" binding:"required,min=2"`
		Customer *string `json:"customer"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id, err := repo.CreateProject(context.Background(), h.DB, req.Name, req.Customer)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
