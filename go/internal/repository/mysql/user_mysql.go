package mysql

import (
	"context"
	"database/sql"

	"ziyadbook/internal/domain"
)

type UserMySQL struct {
	DB *sql.DB
}

func (r UserMySQL) Create(ctx context.Context, u domain.User) (domain.User, error) {
	res, err := r.DB.ExecContext(ctx, "INSERT INTO users(email, name) VALUES(?, ?)", u.Email, u.Name)
	if err != nil {
		return domain.User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}
	u.ID = uint64(id)
	return u, nil
}

func (r UserMySQL) GetByID(ctx context.Context, id uint64) (domain.User, bool, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT id, email, name, created_at FROM users WHERE id = ?", id)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, false, nil
		}
		return domain.User{}, false, err
	}
	return u, true, nil
}

func (r UserMySQL) List(ctx context.Context, limit int) ([]domain.User, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, email, name, created_at FROM users ORDER BY id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
