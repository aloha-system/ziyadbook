package domain

import "time"

type User struct {
	ID        uint64
	Email     string
	Name      string
	CreatedAt time.Time
}
