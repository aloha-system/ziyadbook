package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"ziyadbook/internal/domain"
	"ziyadbook/internal/service"
)

type borrowServiceFake struct {
	borrow domain.Borrow
	err    error
}

func (f *borrowServiceFake) Borrow(ctx context.Context, bookID, memberID uint64) (domain.Borrow, error) {
	return f.borrow, f.err
}

func TestBorrowHandler_Borrow_BadBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := BorrowHandler{Svc: &borrowServiceFake{}}
	h.Register(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/borrow", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBorrowHandler_Borrow_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fake := &borrowServiceFake{
		borrow: domain.Borrow{ID: 1, BookID: 2, MemberID: 3},
	}
	r := gin.New()
	h := BorrowHandler{Svc: fake}
	h.Register(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/borrow", bytes.NewBufferString(`{"book_id":2,"member_id":3}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Contains(t, w.Body.String(), `"book_id":2`)
	require.Contains(t, w.Body.String(), `"member_id":3`)
}

func TestBorrowHandler_Borrow_ErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		serviceErr     error
		expectedStatus int
		expectedBody   string
	}{
		{"book not found", service.ErrBookNotFound, http.StatusNotFound, `"error":"book not found"`},
		{"member not found", service.ErrMemberNotFound, http.StatusNotFound, `"error":"member not found"`},
		{"insufficient stock", service.ErrInsufficientStock, http.StatusConflict, `"error":"insufficient stock"`},
		{"insufficient quota", service.ErrInsufficientQuota, http.StatusConflict, `"error":"insufficient quota"`},
		{"internal", errors.New("boom"), http.StatusInternalServerError, `"error":"internal"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &borrowServiceFake{err: tt.serviceErr}
			r := gin.New()
			h := BorrowHandler{Svc: fake}
			h.Register(r)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/borrow", bytes.NewBufferString(`{"book_id":1,"member_id":1}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			require.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
