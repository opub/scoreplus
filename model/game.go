package model

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/opub/scoreplus/util"
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
	Notes     pq.Int64Array `json:"-"`
	NoteCount int           `sql:"-" json:"notecount"`
}

//Save persists object to data store
func (g *Game) Save() error {
	g.setup()
	if g.ID == 0 {
		g.Created = NullTimeNow()
		return g.execSQL("INSERT INTO game (sport, hometeam, awayteam, homescore, awayscore, start, final, venue, notes, created, createdby) VALUES (:sport, :hometeam, :awayteam, :homescore, :awayscore, :start, :final, :venue, :notes, :created, :createdby) RETURNING id", g)
	}
	g.Modified = NullTimeNow()
	return g.execSQL("UPDATE game SET sport=:sport, hometeam=:hometeam, awayteam=:awayteam, homescore=:homescore, awayscore=:awayscore, start=:start, final=:final, venue=:venue, notes=:notes, modified=:modified, modifiedby=:modifiedby WHERE id=:id", g)
}

//Delete removes object from data store
func (g *Game) Delete() error {
	return g.delete("game")
}

//LinkID gets ID as a linkable string
func (g Game) LinkID() string {
	return util.EncodeLink(g.ID, 10)
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
		g.setup()
		results = append(results, g)
	}
	return results, nil
}

//SelectAllGames from data store
func SelectAllGames() ([]Game, error) {
	return SelectGames(nil)
}

//GetGame returns game from data store
func GetGame(id int64) (Game, error) {
	g := Game{}
	err := get(id, &g)
	g.setup()
	return g, err
}

func (g *Game) setup() {
	g.NoteCount = len(g.Notes)
}
