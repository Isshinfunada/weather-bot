package entity

import "time"

type User struct {
	ID             int
	LINEUserID     string
	SelectedAreaID int
	NotifyTime     time.Time
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
