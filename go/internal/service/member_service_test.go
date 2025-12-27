package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"ziyadbook/internal/domain"
)

type memberRepoFake struct {
	members map[uint64]domain.Member
}

func (f *memberRepoFake) Create(ctx context.Context, m domain.Member) (domain.Member, error) {
	if f.members == nil {
		f.members = make(map[uint64]domain.Member)
	}
	m.ID = uint64(len(f.members) + 1)
	f.members[m.ID] = m
	return m, nil
}

func (f *memberRepoFake) GetByID(ctx context.Context, id uint64) (domain.Member, bool, error) {
	m, ok := f.members[id]
	return m, ok, nil
}

func (f *memberRepoFake) DecrementQuota(ctx context.Context, id uint64, delta uint) (bool, error) {
	m, ok := f.members[id]
	if !ok || m.Quota < delta {
		return false, nil
	}
	m.Quota -= delta
	f.members[id] = m
	return true, nil
}

func (f *memberRepoFake) List(ctx context.Context, limit int) ([]domain.Member, error) {
	out := make([]domain.Member, 0, len(f.members))
	for _, m := range f.members {
		out = append(out, m)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func TestMemberService_Create(t *testing.T) {
	repo := &memberRepoFake{members: map[uint64]domain.Member{}}
	svc := MemberService{Repo: repo}

	m, err := svc.Create(context.Background(), "Member 1", 5)
	require.NoError(t, err)
	require.Equal(t, "Member 1", m.Name)
	require.Equal(t, uint(5), m.Quota)
	require.NotZero(t, m.ID)
}

func TestMemberService_List_Limit(t *testing.T) {
	repo := &memberRepoFake{members: map[uint64]domain.Member{
		1: {ID: 1, Name: "M1", Quota: 1},
		2: {ID: 2, Name: "M2", Quota: 2},
		3: {ID: 3, Name: "M3", Quota: 3},
	}}
	svc := MemberService{Repo: repo}

	members, err := svc.List(context.Background(), 2)
	require.NoError(t, err)
	require.Len(t, members, 2)
}
