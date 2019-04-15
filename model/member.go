package model

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
)

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
	Teams      pq.Int64Array
	Follows    pq.Int64Array
	Followers  pq.Int64Array
}
