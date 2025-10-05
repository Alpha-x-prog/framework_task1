package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/defects/app/backend/internal/repo"
)

type AttachmentsHandler struct{ DB *pgxpool.Pool }

func (h *AttachmentsHandler) Upload(c *gin.Context) {
	defectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file is required"})
		return
	}

	uidVal, _ := c.Get("uid")
	uid, _ := uidVal.(int64)

	// сохранить файл
	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	filename := time.Now().Format("20060102_150405") + "_" + filepath.Base(file.Filename)
	dst := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// простая mime (берём из заголовка)
	mime := file.Header.Get("Content-Type")
	if mime == "" {
		mime = "application/octet-stream"
	}

	id, err := repo.CreateAttachment(c, h.DB, defectID, dst, &mime, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "file_path": dst})
}
