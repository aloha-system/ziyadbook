package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"ziyadbook/internal/domain"
	"ziyadbook/internal/repository"
	"ziyadbook/internal/service"
)

type fakeRepo struct{}

var _ repository.UserRepository = (*fakeRepo)(nil)

func (f *fakeRepo) Create(ctx context.Context, u domain.User) (domain.User, error) {
	u.ID = 1
	return u, nil
}

func (f *fakeRepo) GetByID(ctx context.Context, id uint64) (domain.User, bool, error) {
	return domain.User{}, false, nil
}

func (f *fakeRepo) List(ctx context.Context, limit int) ([]domain.User, error) {
	return []domain.User{}, nil
}

func TestUsersHandler_Create_BadBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := UsersHandler{Svc: service.UserService{Repo: &fakeRepo{}}}
	h.Register(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUsersHandler_Create_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := UsersHandler{Svc: service.UserService{Repo: &fakeRepo{}}}
	h.Register(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("{\"email\":\"a@b.com\",\"name\":\"A\"}"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Contains(t, w.Body.String(), "a@b.com")
}
