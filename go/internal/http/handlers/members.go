package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ziyadbook/internal/domain"
)

// MembersService is the minimal interface MembersHandler depends on.
// Implemented by service.MemberService and test fakes.
type MembersService interface {
	Create(ctx context.Context, name string, quota uint) (domain.Member, error)
	List(ctx context.Context, limit int) ([]domain.Member, error)
}

type MembersHandler struct {
	Svc MembersService
}

type createMemberRequest struct {
	Name  string `json:"name"`
	Quota uint   `json:"quota"`
}

func (h MembersHandler) Register(r gin.IRoutes) {
	r.POST("/members", h.create)
	r.GET("/members", h.list)
}

func (h MembersHandler) create(c *gin.Context) {
	var req createMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, http.StatusBadRequest, "Invalid request body", "ZYD-ERR-310")
		return
	}
	if req.Name == "" {
		WriteError(c, http.StatusBadRequest, "Name is required", "ZYD-ERR-311")
		return
	}

	member, err := h.Svc.Create(c.Request.Context(), req.Name, req.Quota)
	if err != nil {
		WriteError(c, http.StatusInternalServerError, "Terjadi kesalahan internal", "ZYD-ERR-399")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    member.ID,
		"name":  member.Name,
		"quota": member.Quota,
	})
}

func (h MembersHandler) list(c *gin.Context) {
	limit := 20
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	members, err := h.Svc.List(c.Request.Context(), limit)
	if err != nil {
		WriteError(c, http.StatusInternalServerError, "Terjadi kesalahan internal", "ZYD-ERR-301")
		return
	}

	items := make([]gin.H, 0, len(members))
	for _, m := range members {
		items = append(items, gin.H{
			"id":    m.ID,
			"name":  m.Name,
			"quota": m.Quota,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}
