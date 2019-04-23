package model

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
)

//Game data model
type Game struct {
	Base
	Sport     Sport         `json:"sport"`
	HomeTeam  Team          `json:"hometeam"`
	AwayTeam  Team          `json:"awayteam"`
	HomeScore int           `json:"homescore"`
	AwayScore int           `json:"awayscore"`
	Start     null.Time     `json:"start,omitempty"`
	Final     bool          `json:"final"`
	Venue     Venue         `json:"venue,omitempty"`
	Notes     pq.Int64Array `json:"notes,omitempty"`
}

//Save persists object to data store
func (g *Game) Save() error {
	if g.ID == 0 {
		g.Created = nullTimeNow()
		return g.execSQL("INSERT INTO game (sport, hometeam, awayteam, homescore, awayscore, start, final, venue, notes, created, createdby) VALUES (:sport, :hometeam, :awayteam, :homescore, :awayscore, :start, :final, :venue, :notes, :created, :createdby) RETURNING id", g)
	}
	g.Modified = nullTimeNow()
	return g.execSQL("UPDATE game SET sport=:sport, hometeam=:hometeam, awayteam=:awayteam, homescore=:homescore, awayscore=:awayscore, start=:start, final=:final, venue=:venue, notes=:notes, modified=:modified, modifiedby=:modifiedby WHERE id=:id", g)
}

//Delete removes object from data store
func (g *Game) Delete() error {
	return g.delete("game")
}

//SelectGames from data store where ID in slice
func SelectGames(ids []int64) ([]Game, error) {
	rows, err := selectRows(ids, "game")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Game, 0)
	for rows.Next() {
		g := Game{}
		err = rows.StructScan(&g)
		if err != nil {
			return nil, err
		}
		results = append(results, g)
	}
	return results, nil
}

//SelectAllGames from data store
func SelectAllGames() ([]Game, error) {
	return SelectGames(nil)
}

func GetGame(id int64) (Game, error) {
	g := Game{}
	err := Get(id, &g)
	return g, err
}
