package model

import "time"

//Member data model
type Member struct {
	Base
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
}
