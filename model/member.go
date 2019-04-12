package model

import "github.com/guregu/null"

//Member data model
type Member struct {
	Base
	Handle     string
	Email      string
	FirstName  string
	LastName   string
	Verified   bool
	Enabled    bool
	LastActive null.Time
	Teams      []int
	Follows    []int
	Followers  []int
}
