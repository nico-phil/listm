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
	Listnumber  string
	CampaignID  string
	WorkspaceID string
	ListName    string
	Active      bool
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type ListData struct {
	LeadID       string
	WorkspaceID  string
	ListNumber   string
	Firstname    string
	Lastname     string
	PhoneNumber  string
	ZipCode      string
	Extradata    any
	CallCount    int
	CallStatus   string
	Dialable     bool
	InsertedAt   *time.Time
	LastCallDate *time.Time
}
