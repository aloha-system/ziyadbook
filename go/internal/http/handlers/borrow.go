package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ziyadbook/internal/service"
)

type BorrowHandler struct {
	Svc service.BorrowService
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	borrow, err := h.Svc.Borrow(c.Request.Context(), req.BookID, req.MemberID)
	if err != nil {
		switch err {
		case service.ErrBookNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		case service.ErrMemberNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		case service.ErrInsufficientStock:
			c.JSON(http.StatusConflict, gin.H{"error": "insufficient stock"})
		case service.ErrInsufficientQuota:
			c.JSON(http.StatusConflict, gin.H{"error": "insufficient quota"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
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
