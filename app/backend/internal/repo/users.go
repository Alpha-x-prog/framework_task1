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

// Поиск пользователя для логина
func GetUserByEmail(ctx context.Context, db *pgxpool.Pool, email string) (*User, error) {
	row := db.QueryRow(ctx, `
SELECT u.id, u.email, u.password_hash, r.name AS role
FROM users u JOIN roles r ON r.id = u.role_id
WHERE u.email = $1`, email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Hash, &u.Role); err != nil {
		return nil, err
	}
	return &u, nil
}

// Проверка, что email уже занят
func EmailExists(ctx context.Context, db *pgxpool.Pool, email string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}

// Получить id роли по имени (manager/engineer/viewer/lead)
func GetRoleIDByName(ctx context.Context, db *pgxpool.Pool, roleName string) (int, error) {
	var id int
	err := db.QueryRow(ctx, `SELECT id FROM roles WHERE name=$1`, roleName).Scan(&id)
	return id, err
}

// Создать пользователя
func CreateUser(ctx context.Context, db *pgxpool.Pool, email, hash string, roleID int) (int64, error) {
	var id int64
	err := db.QueryRow(ctx,
		`INSERT INTO users(email, password_hash, role_id) VALUES ($1,$2,$3) RETURNING id`,
		email, hash, roleID,
	).Scan(&id)
	return id, err
}

// Публичное представление пользователя (без пароля)
type UserPublic struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// ListUsers — вернуть пользователей с фильтром по роли.
// roleName: "engineer" | "manager" | "viewer" | "lead" | "" (все роли).
// limit/offset — простая пагинация (дефолты и сэйфгварды внутри).
func ListUsers(ctx context.Context, db *pgxpool.Pool, roleName string, limit, offset int) ([]UserPublic, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	const q = `
SELECT u.id, u.email, r.name
FROM users u
JOIN roles r ON r.id = u.role_id
WHERE ($1 = '' OR r.name = $1)
ORDER BY u.email
LIMIT $2 OFFSET $3;
`
	rows, err := db.Query(ctx, q, roleName, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []UserPublic
	for rows.Next() {
		var u UserPublic
		if err := rows.Scan(&u.ID, &u.Email, &u.Role); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// Удобный шорткат: только инженеры
func ListEngineers(ctx context.Context, db *pgxpool.Pool, limit, offset int) ([]UserPublic, error) {
	return ListUsers(ctx, db, "engineer", limit, offset)
}
