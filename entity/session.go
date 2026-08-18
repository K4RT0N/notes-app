package entity

import "time"

type SessionEntity struct {
	TokenHash []byte
	CreatedAt time.Time
	ExpiresAt time.Time
	UserId    int
}
