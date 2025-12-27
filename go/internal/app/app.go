package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"ziyadbook/internal/config"
	apihttp "ziyadbook/internal/http"
	pmysql "ziyadbook/internal/platform/mysql"
	predis "ziyadbook/internal/platform/redis"
	repomysql "ziyadbook/internal/repository/mysql"
	"ziyadbook/internal/service"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	var dbErr error
	var db = (*sql.DB)(nil)
	for i := 0; i < 20; i++ {
		db, dbErr = pmysql.Open(cfg)
		if dbErr == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	var redisErr error
	var redisClient = (*redis.Client)(nil)
	for i := 0; i < 20; i++ {
		redisClient, redisErr = predis.New(cfg)
		if redisErr == nil {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	if redisErr != nil {
		return redisErr
	}
	defer redisClient.Close()

	// Repositories
	userRepo := repomysql.UserMySQL{DB: db}
	bookRepo := repomysql.BookMySQL{DB: db}
	memberRepo := repomysql.MemberMySQL{DB: db}
	borrowRepo := repomysql.BorrowMySQL{DB: db}

	// Services
	userSvc := service.UserService{Repo: userRepo}
	bookSvc := service.BookService{Repo: bookRepo}
	borrowSvc := service.BorrowService{
		BookRepo:   bookRepo,
		MemberRepo: memberRepo,
		BorrowRepo: borrowRepo,
	}
	memberSvc := service.MemberService{Repo: memberRepo}

	r := apihttp.NewRouter(cfg.Env, apihttp.Deps{
		Users:   userSvc,
		Borrows: borrowSvc,
		Books:   bookSvc,
		Members: memberSvc,
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.AppPort),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return srv.ListenAndServe()
}
