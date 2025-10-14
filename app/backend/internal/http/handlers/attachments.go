package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/defects/app/backend/internal/repo"
)

const maxUploadSize = 20 << 20 // 20 MB

type AttachmentsHandler struct {
	DB        *pgxpool.Pool
	UploadDir string // можно прокинуть из конфигурации; если пусто — возьмём из env/дефолт
}

func (h *AttachmentsHandler) getUploadDir() string {
	if h.UploadDir != "" {
		return h.UploadDir
	}
	if v := os.Getenv("UPLOAD_DIR"); v != "" {
		return v
	}

	// дефолт: вложения в отдельной папке
	return filepath.Join("app", "backend", "uploads", "attachments")
}

func (h *AttachmentsHandler) Upload(c *gin.Context) {
	// 1) Валидируем id дефекта
	defectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || defectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}

	// 2) Достаём uid из контекста (обязателен)
	uidVal, ok := c.Get("uid")
	uid, ok2 := uidVal.(int64)
	if !ok || !ok2 || uid <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 3) Ограничиваем размер тела до maxUploadSize (должно быть ДО чтения multipart)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	// 4) Берём файл из формы
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// 5) Готовим директорию
	uploadDir := h.getUploadDir()
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create upload dir"})
		return
	}

	// 6) Безопасное и уникальное имя: YYYYMMDD_HHMMSS_UUID_original.ext
	base := filepath.Base(fileHeader.Filename) // режем любые пути (защита от traversal)
	filename := fmt.Sprintf("%s_%s_%s",
		time.Now().Format("20060102_150405"),
		uuid.NewString(),
		base,
	)
	dstPath := filepath.Join(uploadDir, filename)

	// 7) Сохраняем файл на диск
	if err := c.SaveUploadedFile(fileHeader, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
		return
	}

	// 8) Определяем MIME по содержимому (надёжнее заголовка клиента)
	mime := "application/octet-stream"
	if f, err := os.Open(dstPath); err == nil {
		defer f.Close()
		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		if n > 0 {
			mime = http.DetectContentType(buf[:n])
		}
	}

	// 9) Пишем метаданные в БД (ВАЖНО: используем c.Request.Context())
	id, err := repo.CreateAttachment(c.Request.Context(), h.DB, defectID, dstPath, &mime, uid)
	if err != nil {
		// если запись не удалась — удалим файл, чтобы не копить мусор
		_ = os.Remove(dstPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db insert failed"})
		return
	}

	// 10) Отдаём ответ
	// На проде лучше возвращать публичный URL (например, /uploads/:filename или /api/attachments/:id/download),
	// а не физический путь. Оставляю поле file_path для совместимости с твоим текущим фронтом.
	c.JSON(http.StatusCreated, gin.H{
		"id":        id,
		"file_path": dstPath,
		"mime":      mime,
	})
}

// List attachments for a defect
func (h *AttachmentsHandler) List(c *gin.Context) {
	defectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || defectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	items, err := repo.ListAttachmentsByDefect(c.Request.Context(), h.DB, defectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Download serves a file by attachment id (simple local file serving)
func (h *AttachmentsHandler) Download(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("attId"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	a, err := repo.GetAttachmentByID(c.Request.Context(), h.DB, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	// Best-effort content type
	ct := "application/octet-stream"
	if a.Mime != nil && *a.Mime != "" {
		ct = *a.Mime
	}
	c.Header("Content-Type", ct)
	c.File(a.FilePath)
}
