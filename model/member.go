package model

import "time"

//Member data model
type Member struct {
	ID        string
	Handle    string `sql:",unique"`
	Email     string `sql:",unique"`
	FirstName string
	LastName  string
	Teams     []*Team
	Follows   []*Member
	Followers []*Member
	Created   time.Time `sql:"default:now()"`
	Modified  time.Time
}
