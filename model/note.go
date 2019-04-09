package model

import "time"

//Note data model
type Note struct {
	ID         int
	Message    string
	Created     time.Time `sql:" NOT NULL DEFAULT now()"`
	CreatedBy   int `sql:" NOT NULL"`
	Modified   time.Time
	ModifiedBy int
}
