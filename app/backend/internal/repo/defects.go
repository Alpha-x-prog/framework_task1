package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Defect struct {
	ID          int64      `json:"id"`
	ProjectID   int64      `json:"project_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    int        `json:"priority"`
	AssigneeID  *int64     `json:"assignee_id"`
	StatusID    int        `json:"status_id"`
	DueDate     *time.Time `json:"due_date"`
	CreatedBy   int64      `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type DefectFilter struct {
	ProjectID  *int64
	StatusID   *int
	AssigneeID *int64
	Priority   *int
	Q          *string
	DueFrom    *time.Time
	DueTo      *time.Time
}

func ListDefects(ctx context.Context, db *pgxpool.Pool, f DefectFilter) ([]Defect, error) {
	where := []string{}
	args := []any{}
	i := 1
	add := func(cond string, v any) { where = append(where, fmt.Sprintf(cond, i)); args = append(args, v); i++ }

	if f.ProjectID != nil {
		add("project_id=$%d", *f.ProjectID)
	}
	if f.StatusID != nil {
		add("status_id=$%d", *f.StatusID)
	}
	if f.AssigneeID != nil {
		add("assignee_id=$%d", *f.AssigneeID)
	}
	if f.Priority != nil {
		add("priority=$%d", *f.Priority)
	}
	if f.DueFrom != nil {
		add("due_date>=$%d", *f.DueFrom)
	}
	if f.DueTo != nil {
		add("due_date<=$%d", *f.DueTo)
	}
	if f.Q != nil {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", i, i))
		args = append(args, "%"+*f.Q+"%")
		i++
	}

	sql := `SELECT id, project_id, title, description, priority, assignee_id, status_id, due_date, created_by, created_at, updated_at
	        FROM defects`
	if len(where) > 0 {
		sql += " WHERE " + strings.Join(where, " AND ")
	}
	sql += " ORDER BY id DESC LIMIT 100"

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Defect
	for rows.Next() {
		var d Defect
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.Title, &d.Description, &d.Priority, &d.AssigneeID, &d.StatusID, &d.DueDate, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

type CreateDefectInput struct {
	ProjectID   int64
	Title       string
	Description *string
	Priority    int
	AssigneeID  *int64
	StatusID    int
	DueDate     *time.Time
	CreatedBy   int64
}

func CreateDefect(ctx context.Context, db *pgxpool.Pool, in CreateDefectInput) (int64, error) {
	var id int64
	err := db.QueryRow(ctx, `
INSERT INTO defects(project_id, title, description, priority, assignee_id, status_id, due_date, created_by)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		in.ProjectID, in.Title, in.Description, in.Priority, in.AssigneeID, in.StatusID, in.DueDate, in.CreatedBy).
		Scan(&id)
	return id, err
}

// Доп. функции для статусов (понадобятся на Этапе 7)
func GetDefectByID(ctx context.Context, db *pgxpool.Pool, id int64) (*Defect, error) {
	row := db.QueryRow(ctx, `SELECT id, project_id, title, description, priority, assignee_id, status_id, due_date, created_by, created_at, updated_at FROM defects WHERE id=$1`, id)
	var d Defect
	if err := row.Scan(&d.ID, &d.ProjectID, &d.Title, &d.Description, &d.Priority, &d.AssigneeID, &d.StatusID, &d.DueDate, &d.CreatedBy, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return nil, err
	}
	return &d, nil
}

func UpdateDefectStatus(ctx context.Context, db *pgxpool.Pool, id int64, newStatusID int) error {
	_, err := db.Exec(ctx, `UPDATE defects SET status_id=$1 WHERE id=$2`, newStatusID, id)
	return err
}
