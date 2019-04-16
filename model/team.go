package model

import (
	"github.com/lib/pq"
)

//Team data model
type Team struct {
	Base
	Name   string
	Sport  Sport
	Venue  Venue
	Mascot string
	Games  pq.Int64Array
}

//Save persists object to data store
func (t *Team) Save() error {
	if t.ID == 0 {
		t.Created = nullTimeNow()
		return t.execSQL("INSERT INTO team (name, sport, venue, mascot, games, created, createdby) VALUES (:name, :sport, :venue, :mascot, :games, :created, :createdby) RETURNING id", t)
	}
	t.Modified = nullTimeNow()
	return t.execSQL("UPDATE team SET name=:name, sport=:sport, venue=:venue, mascot=:mascot, games=:games, modified=:modified, modifiedby=:modifiedby WHERE id=:id", t)
}

//Delete removes object from data store
func (t *Team) Delete() error {
	return t.delete("team")
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
		results = append(results, t)
	}
	return results, nil
}

//SelectAllTeams from data store
func SelectAllTeams() ([]Team, error) {
	return SelectTeams(nil)
}
