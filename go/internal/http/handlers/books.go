package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ziyadbook/internal/service"
)

type BooksHandler struct {
	Svc service.BookService
}

func (h BooksHandler) Register(r gin.IRoutes) {
	r.GET("/books", h.list)
	r.POST("/books", h.create)
}

func (h BooksHandler) list(c *gin.Context) {
	limit := 20
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	books, err := h.Svc.List(c.Request.Context(), limit)
	if err != nil {
		WriteError(c, http.StatusInternalServerError, "Terjadi kesalahan internal", "ZYD-ERR-201")
		return
	}

	items := make([]gin.H, 0, len(books))
	for _, b := range books {
		items = append(items, gin.H{
			"id":     b.ID,
			"title":  b.Title,
			"author": b.Author,
			"stock":  b.Stock,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Stock  uint   `json:"stock"`
}

func (h BooksHandler) create(c *gin.Context) {
	var req createBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, http.StatusBadRequest, "Invalid request body", "ZYD-ERR-210")
		return
	}
	if req.Title == "" || req.Author == "" {
		WriteError(c, http.StatusBadRequest, "Title and author are required", "ZYD-ERR-211")
		return
	}

	book, err := h.Svc.Create(c.Request.Context(), req.Title, req.Author, req.Stock)
	if err != nil {
		WriteError(c, http.StatusInternalServerError, "Terjadi kesalahan internal", "ZYD-ERR-299")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":     book.ID,
		"title":  book.Title,
		"author": book.Author,
		"stock":  book.Stock,
	})
}
