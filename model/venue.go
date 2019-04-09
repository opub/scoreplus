package model

import "time"

//Venue data model
type Venue struct {
	ID          int
	Name        string
	Address     string
	Coordinates string
	Created     time.Time `sql:" NOT NULL DEFAULT now()"`
	CreatedBy   int       `sql:" NOT NULL"`
	Modified    time.Time
	ModifiedBy  int
}
