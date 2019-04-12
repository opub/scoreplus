package model

//Sport data model
type Sport struct {
	Base
	Name string
}

//Save persists object to data store
func (b *Sport) Save() error {
	if b.ID == 0 {
		b.Created = NullTimeNow()
		return b.execSQL("INSERT INTO sport (name, created, createdby) VALUES (:name, :created, :createdby) RETURNING id", b)
	}
	b.Modified = NullTimeNow()
	return b.execSQL("UPDATE sport SET name=:name, modified=:modified, modifiedby=:modifiedby WHERE id=:id", b)
}

//Delete removes object from data store
func (b *Sport) Delete() error {
	err := b.execSQL("DELETE FROM sport WHERE id=:id", b)
	if err == nil {
		b.ID = 0
	}
	return err
}
