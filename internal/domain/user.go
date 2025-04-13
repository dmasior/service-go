package domain

import "time"

type User struct {
	ID          string
	Email       string
	LastLoginAt time.Time
	CreatedAt   time.Time
}
