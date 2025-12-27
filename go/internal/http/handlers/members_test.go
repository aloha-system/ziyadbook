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
)

type memberServiceFake struct {
	created domain.Member
	list    []domain.Member
	err     error
}

func (f *memberServiceFake) Create(ctx context.Context, name string, quota uint) (domain.Member, error) {
	return f.created, f.err
}

func (f *memberServiceFake) List(ctx context.Context, limit int) ([]domain.Member, error) {
	return f.list, f.err
}

func TestMembersHandler_Create_BadBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := MembersHandler{Svc: &memberServiceFake{}}
	h.Register(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/members", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	body := w.Body.String()
	require.Contains(t, body, "\"ziyad_error_code\":\"ZYD-ERR-310\"")
}

func TestMembersHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fake := &memberServiceFake{
		created: domain.Member{ID: 1, Name: "M1", Quota: 5},
	}
	r := gin.New()
	h := MembersHandler{Svc: fake}
	h.Register(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/members", bytes.NewBufferString(`{"name":"M1","quota":5}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Contains(t, w.Body.String(), `"name":"M1"`)
}

func TestMembersHandler_List_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fake := &memberServiceFake{
		list: []domain.Member{{ID: 1, Name: "M1", Quota: 5}},
	}
	r := gin.New()
	h := MembersHandler{Svc: fake}
	h.Register(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/members", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"name":"M1"`)
}
