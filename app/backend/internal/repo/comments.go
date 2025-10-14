package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Comment struct {
	ID          int64     `json:"id"`
	DefectID    int64     `json:"defect_id"`
	AuthorID    int64     `json:"author_id"`
	AuthorEmail string    `json:"author_email"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
}

func ListComments(ctx context.Context, db *pgxpool.Pool, defectID int64) ([]Comment, error) {
	rows, err := db.Query(ctx, `
    SELECT c.id, c.defect_id, c.author_id, u.email, c.body, c.created_at
    FROM comments c
    JOIN users u ON u.id = c.author_id
    WHERE c.defect_id=$1
    ORDER BY c.id ASC`, defectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Comment
	for rows.Next() {
		var cmt Comment
		if err := rows.Scan(&cmt.ID, &cmt.DefectID, &cmt.AuthorID, &cmt.AuthorEmail, &cmt.Body, &cmt.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, cmt)
	}
	return out, rows.Err()
}

func CreateComment(ctx context.Context, db *pgxpool.Pool, defectID, authorID int64, body string) (int64, error) {
	var id int64
	err := db.QueryRow(ctx, `INSERT INTO comments(defect_id, author_id, body) VALUES ($1,$2,$3) RETURNING id`, defectID, authorID, body).Scan(&id)
	return id, err
}
