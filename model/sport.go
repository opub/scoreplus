package model

//Sport data model
type Sport struct {
	Base
	Name string
}

//Save persists object to data store
func (b *Sport) Save() error {
	if b.ID == 0 {
		return b.execSQL("INSERT INTO sport (name, createdby, modifiedby) VALUES (:name, 0, 0) RETURNING id", b)
	}
	return b.execSQL("UPDATE sport SET name=:name, modifiedby=1, modified=now() WHERE id=:id", b)
}

//Delete removes object from data store
func (b *Sport) Delete() error {
	err := b.execSQL("DELETE FROM sport WHERE id=:id", b)
	if err == nil {
		b.ID = 0
	}
	return err
}
