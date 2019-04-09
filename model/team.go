package model

import "time"

//Team data model
type Team struct {
	ID         int
	Name       string
	Sport      Sport
	Venue      Venue
	Mascot     string
	Games      []int
	Created     time.Time `sql:" NOT NULL DEFAULT now()"`
	CreatedBy   int `sql:" NOT NULL"`
	Modified   time.Time
	ModifiedBy int
}
