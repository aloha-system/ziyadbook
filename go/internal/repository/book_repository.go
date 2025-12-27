package repository

import (
	"context"

	"ziyadbook/internal/domain"
)

type BookRepository interface {
	Create(ctx context.Context, b domain.Book) (domain.Book, error)
	GetByID(ctx context.Context, id uint64) (domain.Book, bool, error)
	DecrementStock(ctx context.Context, id uint64, delta uint) (bool, error)
	List(ctx context.Context, limit int) ([]domain.Book, error)
}

type MemberRepository interface {
	Create(ctx context.Context, m domain.Member) (domain.Member, error)
	GetByID(ctx context.Context, id uint64) (domain.Member, bool, error)
	DecrementQuota(ctx context.Context, id uint64, delta uint) (bool, error)
	List(ctx context.Context, limit int) ([]domain.Member, error)
}

type BorrowRepository interface {
	Create(ctx context.Context, b domain.Borrow) (domain.Borrow, error)
}
