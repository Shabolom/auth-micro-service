package dto

import "time"

type UserDescription struct {
	AccountID string
	Name      string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
