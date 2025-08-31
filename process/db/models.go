package db

import "time"

// Campaign represents a campaign record from cassandra
type Campaign struct {
	ID            string     `cql:"id"`
	WorkspaceID   string     `cql:"workspace_id"`
	Name          string     `cql:"name"`
	Description   string     `cql:"description"`
	Active        bool       `cql:"active"`
	MaxRatePerMin int        `cql:"max_rate_per_min"`
	DialStartHour int        `cql:"dial_start_hour"`
	DialEndHour   int        `cql:"dial_end_hour"`
	DialDays      []int      `cql:"dial_days"`
	CreatedAt     *time.Time `cql:"createdat"`
	ModifiedAt    *time.Time `cql:"modifiedat"`
}

// List represents a list record from cassandra
type List struct {
	ListNumber  string     `cql:"listnumber"`
	CampaignID  string     `cql:"campaignid"`
	WorkspaceID string     `cql:"workspace_id"`
	ListName    string     `cql:"listname"`
	Active      bool       `cql:"active"`
	CreatedAt   *time.Time `cql:"createdat"`
	UpdatedAt   *time.Time `cql:"updatedat"`
}

// ListData represents a Lead record from cassandra
type ListData struct {
	LeadID       string            `cql:"leadid"`
	ListNumber   string            `cql:"listnumber"`
	WorkspaceID  string            `cql:"workspace_id"`
	PhoneNumber  string            `cql:"phonenumber"`
	FirstName    string            `cql:"firstname"`
	LastName     string            `cql:"lastname"`
	ZipCode      string            `cql:"zipcode"`
	ExtraData    map[string]string `cql:"extradata"`
	CallCount    int               `cql:"callcount"`
	Dialable     bool              `cql:"dialable"`
	InsertedDate time.Time         `cql:"inserteddate"`
	LastCallDate *time.Time        `cql:"lastcalldate"`
	CallStatus   string            `cql:"callstatus"`
}
