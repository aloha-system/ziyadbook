package service

import (
	"context"
	"errors"

	"ziyadbook/internal/domain"
	"ziyadbook/internal/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

var ErrInvalidUser = errors.New("invalid user")

func (s UserService) Create(ctx context.Context, email, name string) (domain.User, error) {
	if email == "" || name == "" {
		return domain.User{}, ErrInvalidUser
	}
	return s.Repo.Create(ctx, domain.User{Email: email, Name: name})
}

func (s UserService) GetByID(ctx context.Context, id uint64) (domain.User, bool, error) {
	if id == 0 {
		return domain.User{}, false, ErrInvalidUser
	}
	return s.Repo.GetByID(ctx, id)
}

func (s UserService) List(ctx context.Context, limit int) ([]domain.User, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.Repo.List(ctx, limit)
}
