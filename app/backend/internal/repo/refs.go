package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ListStatuses(ctx context.Context, db *pgxpool.Pool) ([]RefItem, error) {
	rows, err := db.Query(ctx, `SELECT id, name FROM statuses ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []RefItem
	for rows.Next() {
		var r RefItem
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func ListRoles(ctx context.Context, db *pgxpool.Pool) ([]RefItem, error) {
	rows, err := db.Query(ctx, `SELECT id, name FROM roles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []RefItem
	for rows.Next() {
		var r RefItem
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
