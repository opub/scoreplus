package model

import (
	"database/sql/driver"

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
		t.Created = NullTimeNow()
		return t.execSQL("INSERT INTO team (name, sport, venue, mascot, games, created, createdby) VALUES (:name, :sport, :venue, :mascot, :games, :created, :createdby) RETURNING id", t)
	}
	t.Modified = NullTimeNow()
	return t.execSQL("UPDATE team SET name=:name, sport=:sport, venue=:venue, mascot=:mascot, games=:games, modified=:modified, modifiedby=:modifiedby WHERE id=:id", t)
}

//Delete removes object from data store
func (t *Team) Delete() error {
	err := t.execSQL("DELETE FROM team WHERE id=:id", t)
	if err == nil {
		t.ID = 0
	}
	return err
}

//Scan implements driver Scanner interface
func (t *Team) Scan(value interface{}) error {
	t.ID = value.(int64)
	return nil
}

//Value implements the driver Valuer interface
func (t Team) Value() (driver.Value, error) {
	return t.ID, nil
}
