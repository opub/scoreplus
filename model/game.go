package model

import "github.com/guregu/null"

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
	Notes     []int
}
