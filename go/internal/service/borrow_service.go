package service

import (
	"context"
	"errors"
	"time"

	"ziyadbook/internal/domain"
	"ziyadbook/internal/repository"
)

type BorrowService struct {
	BookRepo   repository.BookRepository
	MemberRepo repository.MemberRepository
	BorrowRepo repository.BorrowRepository
}

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrMemberNotFound    = errors.New("member not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInsufficientQuota = errors.New("insufficient quota")
)

func (s BorrowService) Borrow(ctx context.Context, bookID, memberID uint64) (domain.Borrow, error) {
	// Validate existence
	_, okBook, err := s.BookRepo.GetByID(ctx, bookID)
	if err != nil {
		return domain.Borrow{}, err
	}
	if !okBook {
		return domain.Borrow{}, ErrBookNotFound
	}

	_, okMember, err := s.MemberRepo.GetByID(ctx, memberID)
	if err != nil {
		return domain.Borrow{}, err
	}
	if !okMember {
		return domain.Borrow{}, ErrMemberNotFound
	}

	// Perform atomic decrement checks within a transaction
	// For simplicity, we rely on repository-level atomic updates
	stockOK, err := s.BookRepo.DecrementStock(ctx, bookID, 1)
	if err != nil {
		return domain.Borrow{}, err
	}
	if !stockOK {
		return domain.Borrow{}, ErrInsufficientStock
	}

	quotaOK, err := s.MemberRepo.DecrementQuota(ctx, memberID, 1)
	if err != nil {
		return domain.Borrow{}, err
	}
	if !quotaOK {
		return domain.Borrow{}, ErrInsufficientQuota
	}

	borrow := domain.Borrow{
		BookID:     bookID,
		MemberID:   memberID,
		BorrowedAt: time.Now(),
	}
	created, err := s.BorrowRepo.Create(ctx, borrow)
	if err != nil {
		return domain.Borrow{}, err
	}
	return created, nil
}
