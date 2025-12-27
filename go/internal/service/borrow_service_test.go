package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"ziyadbook/internal/domain"
)

type fakeBookRepo struct {
	books map[uint64]domain.Book
}

func (f *fakeBookRepo) Create(ctx context.Context, b domain.Book) (domain.Book, error) {
	b.ID = uint64(len(f.books) + 1)
	f.books[b.ID] = b
	return b, nil
}

func (f *fakeBookRepo) GetByID(ctx context.Context, id uint64) (domain.Book, bool, error) {
	b, ok := f.books[id]
	return b, ok, nil
}

func (f *fakeBookRepo) DecrementStock(ctx context.Context, id uint64, delta uint) (bool, error) {
	b, ok := f.books[id]
	if !ok || b.Stock < delta {
		return false, nil
	}
	b.Stock -= delta
	f.books[id] = b
	return true, nil
}

type fakeMemberRepo struct {
	members map[uint64]domain.Member
}

func (f *fakeMemberRepo) Create(ctx context.Context, m domain.Member) (domain.Member, error) {
	m.ID = uint64(len(f.members) + 1)
	f.members[m.ID] = m
	return m, nil
}

func (f *fakeMemberRepo) GetByID(ctx context.Context, id uint64) (domain.Member, bool, error) {
	m, ok := f.members[id]
	return m, ok, nil
}

func (f *fakeMemberRepo) DecrementQuota(ctx context.Context, id uint64, delta uint) (bool, error) {
	m, ok := f.members[id]
	if !ok || m.Quota < delta {
		return false, nil
	}
	m.Quota -= delta
	f.members[id] = m
	return true, nil
}

type fakeBorrowRepo struct {
	borrows []domain.Borrow
}

func (f *fakeBorrowRepo) Create(ctx context.Context, b domain.Borrow) (domain.Borrow, error) {
	b.ID = uint64(len(f.borrows) + 1)
	f.borrows = append(f.borrows, b)
	return b, nil
}

func TestBorrowService_Borrow_Success(t *testing.T) {
	bookRepo := &fakeBookRepo{books: map[uint64]domain.Book{
		1: {ID: 1, Title: "Go", Author: "A", Stock: 3},
	}}
	memberRepo := &fakeMemberRepo{members: map[uint64]domain.Member{
		1: {ID: 1, Name: "M", Quota: 2},
	}}
	borrowRepo := &fakeBorrowRepo{}

	svc := BorrowService{
		BookRepo:   bookRepo,
		MemberRepo: memberRepo,
		BorrowRepo: borrowRepo,
	}

	borrow, err := svc.Borrow(context.Background(), 1, 1)
	require.NoError(t, err)
	require.Equal(t, uint64(1), borrow.BookID)
	require.Equal(t, uint64(1), borrow.MemberID)
	require.NotZero(t, borrow.ID)

	book, ok := bookRepo.books[1]
	require.True(t, ok)
	require.Equal(t, uint(2), book.Stock)

	member, ok := memberRepo.members[1]
	require.True(t, ok)
	require.Equal(t, uint(1), member.Quota)
}

func TestBorrowService_Borrow_Errors(t *testing.T) {
	bookRepo := &fakeBookRepo{books: map[uint64]domain.Book{
		1: {ID: 1, Title: "Go", Author: "A", Stock: 0},
	}}
	memberRepo := &fakeMemberRepo{members: map[uint64]domain.Member{
		1: {ID: 1, Name: "M", Quota: 0},
	}}
	borrowRepo := &fakeBorrowRepo{}

	svc := BorrowService{
		BookRepo:   bookRepo,
		MemberRepo: memberRepo,
		BorrowRepo: borrowRepo,
	}

	_, err := svc.Borrow(context.Background(), 99, 1)
	require.ErrorIs(t, err, ErrBookNotFound)

	_, err = svc.Borrow(context.Background(), 1, 99)
	require.ErrorIs(t, err, ErrMemberNotFound)

	_, err = svc.Borrow(context.Background(), 1, 1)
	require.ErrorIs(t, err, ErrInsufficientStock)

	// Give stock but no quota
	bookRepo.books[1] = domain.Book{ID: 1, Title: "Go", Author: "A", Stock: 1}
	_, err = svc.Borrow(context.Background(), 1, 1)
	require.ErrorIs(t, err, ErrInsufficientQuota)
}
