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

func ListAttachmentsByDefect(ctx context.Context, db *pgxpool.Pool, defectID int64) ([]Attachment, error) {
	rows, err := db.Query(ctx, `SELECT id, defect_id, file_path, mime, uploaded_by, created_at FROM attachments WHERE defect_id=$1 ORDER BY id DESC`, defectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Attachment
	for rows.Next() {
		var a Attachment
		if err := rows.Scan(&a.ID, &a.DefectID, &a.FilePath, &a.Mime, &a.UploadedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func GetAttachmentByID(ctx context.Context, db *pgxpool.Pool, id int64) (*Attachment, error) {
	row := db.QueryRow(ctx, `SELECT id, defect_id, file_path, mime, uploaded_by, created_at FROM attachments WHERE id=$1`, id)
	var a Attachment
	if err := row.Scan(&a.ID, &a.DefectID, &a.FilePath, &a.Mime, &a.UploadedBy, &a.CreatedAt); err != nil {
		return nil, err
	}
	return &a, nil
}
