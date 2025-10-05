package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Attachment struct {
	ID         int64     `json:"id"`
	DefectID   int64     `json:"defect_id"`
	FilePath   string    `json:"file_path"`
	Mime       *string   `json:"mime"`
	UploadedBy int64     `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}

func CreateAttachment(ctx context.Context, db *pgxpool.Pool, defectID int64, filePath string, mime *string, uploadedBy int64) (int64, error) {
	var id int64
	err := db.QueryRow(ctx, `INSERT INTO attachments(defect_id, file_path, mime, uploaded_by) VALUES ($1,$2,$3,$4) RETURNING id`,
		defectID, filePath, mime, uploadedBy).Scan(&id)
	return id, err
}
