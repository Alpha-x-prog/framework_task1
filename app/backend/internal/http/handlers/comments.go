package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/defects/app/backend/internal/repo"
)

type CommentsHandler struct{ DB *pgxpool.Pool }

func (h *CommentsHandler) List(c *gin.Context) {
	defectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad id"})
		return
	}
	items, err := repo.ListComments(c, h.DB, defectID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *CommentsHandler) Create(c *gin.Context) {
	defectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad id"})
		return
	}
	var req struct {
		Body string `json:"body" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	uidVal, _ := c.Get("uid")
	uid, _ := uidVal.(int64)

	id, err := repo.CreateComment(c, h.DB, defectID, uid, req.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
