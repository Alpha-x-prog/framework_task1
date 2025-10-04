package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Project struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Customer *string `json:"customer"`
}

func ListProjects(ctx context.Context, db *pgxpool.Pool) ([]Project, error) {
	rows, err := db.Query(ctx, `SELECT id, name, customer FROM projects ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Customer); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func CreateProject(ctx context.Context, db *pgxpool.Pool, name string, customer *string) (int64, error) {
	var id int64
	err := db.QueryRow(ctx,
		`INSERT INTO projects(name, customer) VALUES ($1, $2) RETURNING id`,
		name, customer).Scan(&id)
	return id, err
}
