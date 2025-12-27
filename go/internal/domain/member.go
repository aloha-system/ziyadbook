package domain

import "time"

type Member struct {
	ID        uint64
	Name      string
	Quota     uint
	CreatedAt time.Time
}
