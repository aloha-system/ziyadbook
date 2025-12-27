package http

import (
	"github.com/gin-gonic/gin"

	"ziyadbook/internal/http/handlers"
	"ziyadbook/internal/service"
)

type Deps struct {
	Users   service.UserService
	Borrows service.BorrowService
	Books   service.BookService
	Members service.MemberService
}

func NewRouter(env string, deps Deps) *gin.Engine {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/")
	(handlers.HealthHandler{}).Register(api)
	(handlers.UsersHandler{Svc: deps.Users}).Register(api)
	(handlers.BorrowHandler{Svc: deps.Borrows}).Register(api)
	(handlers.BooksHandler{Svc: deps.Books}).Register(api)
	(handlers.MembersHandler{Svc: deps.Members}).Register(api)

	return r
}
