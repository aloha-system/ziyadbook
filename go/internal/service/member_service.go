package service

import (
	"context"

	"ziyadbook/internal/domain"
	"ziyadbook/internal/repository"
)

type MemberService struct {
	Repo repository.MemberRepository
}

func (s MemberService) Create(ctx context.Context, name string, quota uint) (domain.Member, error) {
	m := domain.Member{
		Name:  name,
		Quota: quota,
	}
	return s.Repo.Create(ctx, m)
}

func (s MemberService) List(ctx context.Context, limit int) ([]domain.Member, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.Repo.List(ctx, limit)
}
