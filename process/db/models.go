package db

import "time"

type Campaign struct {
	ID            string
	WorkspaceID   string
	Name          string
	Description   string
	Active        bool
	MaxRatePerMin int
	DialStartHour int
	DialEndHour   int
	DialDays      []int
	CreatedAt     *time.Time
	ModifiedAt    *time.Time
}

type List struct {
}

type ListData struct {
}
