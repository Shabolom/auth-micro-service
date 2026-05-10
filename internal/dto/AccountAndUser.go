package dto

import "time"

type AccountAndUser struct {
	ID    string
	Email string

	Name *string
	Age  *int

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
