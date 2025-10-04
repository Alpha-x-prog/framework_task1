// app/backend/internal/repo/users.go
package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID    int64
	Email string
	Hash  string
	Role  string
}

func GetUserByEmail(ctx context.Context, db *pgxpool.Pool, email string) (*User, error) {
	q := `
SELECT u.id, u.email, u.password_hash, r.name AS role
FROM users u
JOIN roles r ON r.id = u.role_id
WHERE u.email = $1`
	row := db.QueryRow(ctx, q, email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Hash, &u.Role); err != nil {
		return nil, err
	}
	return &u, nil
}
