package model

import "time"

//Game data model
type Game struct {
	ID         int
	Sport      Sport
	HomeTeam   Team
	AwayTeam   Team
	HomeScore  int
	AwayScore  int
	Start      time.Time
	Final      bool
	Venue      Venue
	Date       time.Time
	Notes      []int
	Created    time.Time `sql:" NOT NULL DEFAULT now()"`
	CreatedBy  int       `sql:" NOT NULL"`
	Modified   time.Time
	ModifiedBy int
}
