package model

import "github.com/lib/pq"

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
func (b *Team) Save() error {
	if b.ID == 0 {
		b.Created = NullTimeNow()
		return b.execSQL("INSERT INTO team (name, sport, venue, mascot, games, created, createdby) VALUES (:name, :sport, :venue, :mascot, :games, :created, :createdby) RETURNING id", b)
	}
	b.Modified = NullTimeNow()
	return b.execSQL("UPDATE team SET name=:name, sport=:sport, venue=:venue, mascot=:mascot, games=:games, modified=:modified, modifiedby=:modifiedby WHERE id=:id", b)
}

//Delete removes object from data store
func (b *Team) Delete() error {
	err := b.execSQL("DELETE FROM team WHERE id=:id", b)
	if err == nil {
		b.ID = 0
	}
	return err
}
