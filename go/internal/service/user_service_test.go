package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"ziyadbook/internal/domain"
)

type fakeUserRepo struct {
	created []domain.User
	get     map[uint64]domain.User
}

func (f *fakeUserRepo) Create(ctx context.Context, u domain.User) (domain.User, error) {
	u.ID = uint64(len(f.created) + 1)
	f.created = append(f.created, u)
	if f.get == nil {
		f.get = map[uint64]domain.User{}
	}
	f.get[u.ID] = u
	return u, nil
}

func (f *fakeUserRepo) GetByID(ctx context.Context, id uint64) (domain.User, bool, error) {
	u, ok := f.get[id]
	return u, ok, nil
}

func (f *fakeUserRepo) List(ctx context.Context, limit int) ([]domain.User, error) {
	out := make([]domain.User, 0)
	for i := len(f.created) - 1; i >= 0; i-- {
		out = append(out, f.created[i])
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func TestUserService_Create_Validates(t *testing.T) {
	s := UserService{Repo: &fakeUserRepo{}}
	_, err := s.Create(context.Background(), "", "x")
	require.ErrorIs(t, err, ErrInvalidUser)
}

func TestUserService_Create_Success(t *testing.T) {
	s := UserService{Repo: &fakeUserRepo{}}
	u, err := s.Create(context.Background(), "a@b.com", "A")
	require.NoError(t, err)
	require.NotZero(t, u.ID)
	require.Equal(t, "a@b.com", u.Email)
}
