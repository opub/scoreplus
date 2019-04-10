package model

import "time"

//Game data model
type Game struct {
	Base
	Sport     Sport
	HomeTeam  Team
	AwayTeam  Team
	HomeScore int
	AwayScore int
	Start     time.Time
	Final     bool
	Venue     Venue
	Notes     []int
}
