package model

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
)

//Game data model
type Game struct {
	Base
	Sport     Sport
	HomeTeam  Team
	AwayTeam  Team
	HomeScore int
	AwayScore int
	Start     null.Time
	Final     bool
	Venue     Venue
	Notes     pq.Int64Array
}
