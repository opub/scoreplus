package model

//Venue data model
type Venue struct {
	Base
	Name        string
	Address     string
	Coordinates string
}

//Save persists object to data store
func (b *Venue) Save() error {
	if b.ID == 0 {
		b.Created = NullTimeNow()
		return b.execSQL("INSERT INTO venue (name, address, coordinates, created, createdby) VALUES (:name, :address, :coordinates, :created, :createdby) RETURNING id", b)
	}
	b.Modified = NullTimeNow()
	return b.execSQL("UPDATE venue SET name=:name, address=:address, coordinates=:coordinates, modified=:modified, modifiedby=:modifiedby WHERE id=:id", b)
}

//Delete removes object from data store
func (b *Venue) Delete() error {
	err := b.execSQL("DELETE FROM venue WHERE id=:id", b)
	if err == nil {
		b.ID = 0
	}
	return err
}
