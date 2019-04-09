package model

import "time"

//Game data model
type Game struct {
	ID        string
	Sport     Sport
	HomeTeam  Team
	AwayTeam  Team
	HomeScore int
	AwayScore int
	Final     bool
	Location  string
	Date      time.Time
	Notes     []string
	Created   time.Time `sql:"default:now()"`
	Modified  time.Time
}
