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

//Save persists object to data store
func (g *Game) Save() error {
	if g.ID == 0 {
		g.Created = NullTimeNow()
		return g.execSQL("INSERT INTO game (sport, hometeam, awayteam, homescore, awayscore, start, final, venue, notes, created, createdby) VALUES (:sport, :hometeam, :awayteam, :homescore, :awayscore, :start, :final, :venue, :notes, :created, :createdby) RETURNING id", g)
	}
	g.Modified = NullTimeNow()
	return g.execSQL("UPDATE game SET sport=:sport, hometeam=:hometeam, awayteam=:awayteam, homescore=:homescore, awayscore=:awayscore, start=:start, final=:final, venue=:venue, notes=:notes, modified=:modified, modifiedby=:modifiedby WHERE id=:id", g)
}

//Delete removes object from data store
func (g *Game) Delete() error {
	err := g.execSQL("DELETE FROM game WHERE id=:id", g)
	if err == nil {
		g.ID = 0
	}
	return err
}
