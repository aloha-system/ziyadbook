package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"ziyadbook/internal/domain"
	"ziyadbook/internal/service"
)

// BorrowService is the minimal interface BorrowHandler depends on.
// It is implemented by service.BorrowService and by test fakes.
type BorrowService interface {
	Borrow(ctx context.Context, bookID, memberID uint64) (domain.Borrow, error)
}

type BorrowHandler struct {
	Svc BorrowService
}

type borrowRequest struct {
	BookID   uint64 `json:"book_id"`
	MemberID uint64 `json:"member_id"`
}

func (h BorrowHandler) Register(r gin.IRoutes) {
	r.POST("/borrow", h.borrow)
}

func (h BorrowHandler) borrow(c *gin.Context) {
	var req borrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, http.StatusBadRequest, "Invalid request body", "ZYD-ERR-000")
		return
	}

	borrow, err := h.Svc.Borrow(c.Request.Context(), req.BookID, req.MemberID)
	if err != nil {
		switch err {
		case service.ErrBookNotFound:
			WriteError(c, http.StatusNotFound, "Buku tidak ditemukan", "ZYD-ERR-002")
		case service.ErrMemberNotFound:
			WriteError(c, http.StatusNotFound, "Member tidak ditemukan", "ZYD-ERR-003")
		case service.ErrInsufficientStock:
			WriteError(c, http.StatusConflict, "Stok buku habis", "ZYD-ERR-001")
		case service.ErrInsufficientQuota:
			WriteError(c, http.StatusConflict, "Kuota peminjaman member habis", "ZYD-ERR-004")
		default:
			WriteError(c, http.StatusInternalServerError, "Terjadi kesalahan internal", "ZYD-ERR-999")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          borrow.ID,
		"book_id":     borrow.BookID,
		"member_id":   borrow.MemberID,
		"borrowed_at": borrow.BorrowedAt,
	})
}
