package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"ziyadbook/internal/config"
)

func Open(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true",
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDatabase,
	)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := DB.PingContext(ctx); err != nil {
		_ = DB.Close()
		return nil, err
	}

	return DB, nil
}
