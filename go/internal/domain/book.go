package domain

import "time"

type Book struct {
	ID        uint64
	Title     string
	Author    string
	Stock     uint
	CreatedAt time.Time
}
