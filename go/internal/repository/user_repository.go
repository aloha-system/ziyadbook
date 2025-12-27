package repository

import (
	"context"

	"ziyadbook/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	GetByID(ctx context.Context, id uint64) (domain.User, bool, error)
	List(ctx context.Context, limit int) ([]domain.User, error)
}
