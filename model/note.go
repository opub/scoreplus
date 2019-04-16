package model

//Note data model
type Note struct {
	Base
	Message string
}

//Save persists object to data store
func (n *Note) Save() error {
	if n.ID == 0 {
		n.Created = nullTimeNow()
		return n.execSQL("INSERT INTO note (message, created, createdby) VALUES (:message, :created, :createdby) RETURNING id", n)
	}
	n.Modified = nullTimeNow()
	return n.execSQL("UPDATE note SET message=:message, modified=:modified, modifiedby=:modifiedby WHERE id=:id", n)
}

//Delete removes object from data store
func (n *Note) Delete() error {
	return n.delete("note")
}
