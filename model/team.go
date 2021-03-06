package model

import (
	"strings"

	"github.com/lib/pq"
	"github.com/opub/scoreplus/db"
	"github.com/opub/scoreplus/util"
)

//Team data model
type Team struct {
	Base
	Name      string        `json:"name"`
	Sport     Sport         `json:"sport,omitempty"`
	Venue     Venue         `json:"venue,omitempty"`
	Mascot    string        `json:"mascot,omitempty"`
	Games     pq.Int64Array `json:"-"`
	GameCount int           `sql:"-" json:"gameCount"`
}

//Save persists object to data store
func (t *Team) Save() error {
	t.setup()
	if t.ID == 0 {
		t.Created = NullTimeNow()
		return t.execSQL("INSERT INTO team (name, sport, venue, mascot, games, created, createdby) VALUES (:name, :sport, :venue, :mascot, :games, :created, :createdby) RETURNING id", t)
	}
	t.Modified = NullTimeNow()
	return t.execSQL("UPDATE team SET name=:name, sport=:sport, venue=:venue, mascot=:mascot, games=:games, modified=:modified, modifiedby=:modifiedby WHERE id=:id", t)
}

//Delete removes object from data store
func (t *Team) Delete() error {
	return t.delete("team")
}

//LinkID gets ID as a linkable string
func (t Team) LinkID() string {
	return util.EncodeLink(t.ID, 40)
}

//SelectTeams from data store where ID in slice
func SelectTeams(ids []int64) ([]Team, error) {
	rows, err := selectRows(ids, "team")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Team, 0)
	for rows.Next() {
		t := Team{}
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}
		t.setup()
		results = append(results, t)
	}
	return results, nil
}

//SearchTeams finds teams that match a search string and sport
func SearchTeams(search string, sport string) ([]Team, error) {
	db, err := db.Connect()
	if err != nil {
		return nil, err
	}

	term := strings.ToLower(strings.TrimSpace(search))
	sql := "SELECT * FROM team WHERE lower(name) LIKE '%' || $1 || '%' AND sport = $2"
	rows, err := db.Queryx(sql, term, sport)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Team, 0)
	for rows.Next() {
		t := Team{}
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, nil
}

//SelectAllTeams from data store
func SelectAllTeams() ([]Team, error) {
	return SelectTeams(nil)
}

//GetTeam returns team from data store
func GetTeam(id int64) (Team, error) {
	t := Team{}
	err := get(id, &t)
	t.setup()
	return t, err
}

//GetTeamFull returns team from data store with nested structs populated
func GetTeamFull(id int64) (Team, error) {
	t, err := GetTeam(id)
	if t.Venue.ID != 0 && err == nil {
		t.Venue, err = GetVenue(t.Venue.ID)
	}
	return t, err
}

func (t *Team) setup() {
	t.GameCount = len(t.Games)
}
