package handlers

import (
	"net/http"
	"strconv"
	"time"

	"example/defects/app/backend/internal/core"
	"example/defects/app/backend/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DefectsHandler struct{ DB *pgxpool.Pool }

func asIntPtr(s string) *int {
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}
func asInt64Ptr(s string) *int64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &v
}
func asDatePtr(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

func (h *DefectsHandler) List(c *gin.Context) {
	f := repo.DefectFilter{
		ID:         asInt64Ptr(c.Query("id")),
		ProjectID:  asInt64Ptr(c.Query("project_id")),
		StatusID:   asIntPtr(c.Query("status_id")),
		AssigneeID: asInt64Ptr(c.Query("assignee_id")),
		Priority:   asIntPtr(c.Query("priority")),
		DueFrom:    asDatePtr(c.Query("due_from")),
		DueTo:      asDatePtr(c.Query("due_to")),
	}
	if q := c.Query("q"); q != "" {
		f.Q = &q
	}

	items, err := repo.ListDefects(c, h.DB, f)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *DefectsHandler) Create(c *gin.Context) {
	var req struct {
		ProjectID   int64   `json:"project_id" binding:"required"`
		Title       string  `json:"title" binding:"required,min=3"`
		Description *string `json:"description"`
		Priority    *int    `json:"priority"` // default 3
		AssigneeID  *int64  `json:"assignee_id"`
		StatusID    *int    `json:"status_id"` // id статуса (например, new)
		DueDate     *string `json:"due_date"`  // "YYYY-MM-DD"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	uidVal, _ := c.Get("uid")
	uid, _ := uidVal.(int64)

	priority := 3
	if req.Priority != nil {
		priority = *req.Priority
	}
	statusID := 0
	if req.StatusID != nil {
		statusID = *req.StatusID
	}

	var due *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		if t, err := time.Parse("2006-01-02", *req.DueDate); err == nil {
			due = &t
		}
	}

	id, err := repo.CreateDefect(c, h.DB, repo.CreateDefectInput{
		ProjectID: req.ProjectID, Title: req.Title, Description: req.Description,
		Priority: priority, AssigneeID: req.AssigneeID, StatusID: statusID,
		DueDate: due, CreatedBy: uid,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *DefectsHandler) UpdateStatus(c *gin.Context) {
	// путь: /api/defects/:id/status
	defectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad id"})
		return
	}

	var req struct {
		StatusID int `json:"status_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// получить текущий дефект
	cur, err := repo.GetDefectByID(c, h.DB, defectID)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	if !core.CanTransit(cur.StatusID, req.StatusID) {
		c.JSON(400, gin.H{"error": "invalid transition"})
		return
	}

	if err := repo.UpdateDefectStatus(c, h.DB, defectID, req.StatusID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
