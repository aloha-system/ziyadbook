package mysql

import (
	"context"
	"database/sql"

	"ziyadbook/internal/domain"
)

type BorrowMySQL struct {
	DB *sql.DB
}

func (r BorrowMySQL) Create(ctx context.Context, b domain.Borrow) (domain.Borrow, error) {
	res, err := r.DB.ExecContext(ctx, "INSERT INTO borrows(book_id, member_id, borrowed_at) VALUES(?, ?, ?)", b.BookID, b.MemberID, b.BorrowedAt)
	if err != nil {
		return domain.Borrow{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Borrow{}, err
	}
	b.ID = uint64(id)
	return b, nil
}
