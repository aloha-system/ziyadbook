package domain

import "time"

type Borrow struct {
ID         uint64
BookID     uint64
MemberID   uint64
BorrowedAt time.Time
ReturnedAt *time.Time
}
