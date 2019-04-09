package model

import "time"

//Member data model
type Member struct {
	ID         int
	Handle     string
	Email      string
	FirstName  string
	LastName   string
	Verified   bool
	Enabled    bool
	LastActive time.Time
	Teams      []int
	Follows    []int
	Followers  []int
	Created    time.Time `sql:" NOT NULL DEFAULT now()"`
	CreatedBy  int       `sql:" NOT NULL"`
	Modified   time.Time
	ModifiedBy int
}
