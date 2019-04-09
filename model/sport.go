package model

import "time"

//Sport data model
type Sport struct {
	ID       string
	Name     string    `sql:",unique"`
	Created  time.Time `sql:"default:now()"`
	Modified time.Time
}
