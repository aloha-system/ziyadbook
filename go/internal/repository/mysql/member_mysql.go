package mysql

import (
	"context"
	"database/sql"

	"ziyadbook/internal/domain"
)

type MemberMySQL struct {
	DB *sql.DB
}

func (r MemberMySQL) Create(ctx context.Context, m domain.Member) (domain.Member, error) {
	res, err := r.DB.ExecContext(ctx, "INSERT INTO members(name, quota) VALUES(?, ?)", m.Name, m.Quota)
	if err != nil {
		return domain.Member{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Member{}, err
	}
	m.ID = uint64(id)
	return m, nil
}

func (r MemberMySQL) GetByID(ctx context.Context, id uint64) (domain.Member, bool, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT id, name, quota, created_at FROM members WHERE id = ?", id)
	var m domain.Member
	if err := row.Scan(&m.ID, &m.Name, &m.Quota, &m.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Member{}, false, nil
		}
		return domain.Member{}, false, err
	}
	return m, true, nil
}

func (r MemberMySQL) DecrementQuota(ctx context.Context, id uint64, delta uint) (bool, error) {
	res, err := r.DB.ExecContext(ctx, "UPDATE members SET quota = quota - ? WHERE id = ? AND quota >= ?", delta, id, delta)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (r MemberMySQL) List(ctx context.Context, limit int) ([]domain.Member, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, name, quota, created_at FROM members ORDER BY id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.Member
	for rows.Next() {
		var m domain.Member
		if err := rows.Scan(&m.ID, &m.Name, &m.Quota, &m.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}
