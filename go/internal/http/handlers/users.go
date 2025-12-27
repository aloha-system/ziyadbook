package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ziyadbook/internal/service"
)

type UsersHandler struct {
	Svc service.UserService
}

type createUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h UsersHandler) Register(r gin.IRoutes) {
	r.POST("/users", h.create)
	r.GET("/users/:id", h.getByID)
	r.GET("/users", h.list)
}

func (h UsersHandler) create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	u, err := h.Svc.Create(c.Request.Context(), req.Email, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email, "name": u.Name})
}

func (h UsersHandler) getByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	u, ok, err := h.Svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email, "name": u.Name})
}

func (h UsersHandler) list(c *gin.Context) {
	limit := 20
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	users, err := h.Svc.List(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	out := make([]gin.H, 0, len(users))
	for _, u := range users {
		out = append(out, gin.H{"id": u.ID, "email": u.Email, "name": u.Name})
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}
