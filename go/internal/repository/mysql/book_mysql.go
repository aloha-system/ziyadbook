package mysql

import (
	"context"
	"database/sql"

	"ziyadbook/internal/domain"
)

type BookMySQL struct {
	DB *sql.DB
}

func (r BookMySQL) Create(ctx context.Context, b domain.Book) (domain.Book, error) {
	res, err := r.DB.ExecContext(ctx, "INSERT INTO books(title, author, stock) VALUES(?, ?, ?)", b.Title, b.Author, b.Stock)
	if err != nil {
		return domain.Book{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Book{}, err
	}
	b.ID = uint64(id)
	return b, nil
}

func (r BookMySQL) GetByID(ctx context.Context, id uint64) (domain.Book, bool, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT id, title, author, stock, created_at FROM books WHERE id = ?", id)
	var b domain.Book
	if err := row.Scan(&b.ID, &b.Title, &b.Author, &b.Stock, &b.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Book{}, false, nil
		}
		return domain.Book{}, false, err
	}
	return b, true, nil
}

func (r BookMySQL) DecrementStock(ctx context.Context, id uint64, delta uint) (bool, error) {
	res, err := r.DB.ExecContext(ctx, "UPDATE books SET stock = stock - ? WHERE id = ? AND stock >= ?", delta, id, delta)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (r BookMySQL) List(ctx context.Context, limit int) ([]domain.Book, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, title, author, stock, created_at FROM books ORDER BY id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Stock, &b.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}
