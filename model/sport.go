package model

import "time"

//Sport data model
type Sport struct {
	ID         int
	Name       string
	Created     time.Time `sql:" NOT NULL DEFAULT now()"`
	CreatedBy   int `sql:" NOT NULL"`
	Modified   time.Time
	ModifiedBy int
}
