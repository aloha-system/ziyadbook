package service

import (
	"context"

	"ziyadbook/internal/domain"
	"ziyadbook/internal/repository"
)

type BookService struct {
	Repo repository.BookRepository
}

func (s BookService) List(ctx context.Context, limit int) ([]domain.Book, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.Repo.List(ctx, limit)
}

func (s BookService) Create(ctx context.Context, title, author string, stock uint) (domain.Book, error) {
	book := domain.Book{
		Title:  title,
		Author: author,
		Stock:  stock,
	}
	return s.Repo.Create(ctx, book)
}
