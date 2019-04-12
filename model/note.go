package model

//Note data model
type Note struct {
	Base
	Message string
}

//Save persists object to data store
func (b *Note) Save() error {
	if b.ID == 0 {
		b.Created = NullTimeNow()
		return b.execSQL("INSERT INTO note (message, created, createdby) VALUES (:message, :created, :createdby) RETURNING id", b)
	}
	b.Modified = NullTimeNow()
	return b.execSQL("UPDATE note SET message=:message, modified=:modified, modifiedby=:modifiedby WHERE id=:id", b)
}

//Delete removes object from data store
func (b *Note) Delete() error {
	err := b.execSQL("DELETE FROM note WHERE id=:id", b)
	if err == nil {
		b.ID = 0
	}
	return err
}
