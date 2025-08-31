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
	ListNumber  string     `cql:"listnumber"`
	CampaignID  string     `cql:"campaignid"`
	WorkspaceID string     `cql:"workspace_id"`
	ListName    string     `cql:"listname"`
	Active      bool       `cql:"active"`
	CreatedAt   *time.Time `cql:"createdat"`
	UpdatedAt   *time.Time `cql:"updatedat"`
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
