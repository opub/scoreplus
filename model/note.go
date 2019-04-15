package model

import "database/sql/driver"

//Note data model
type Note struct {
	Base
	Message string
}

//Save persists object to data store
func (n *Note) Save() error {
	if n.ID == 0 {
		n.Created = NullTimeNow()
		return n.execSQL("INSERT INTO note (message, created, createdby) VALUES (:message, :created, :createdby) RETURNING id", n)
	}
	n.Modified = NullTimeNow()
	return n.execSQL("UPDATE note SET message=:message, modified=:modified, modifiedby=:modifiedby WHERE id=:id", n)
}

//Delete removes object from data store
func (n *Note) Delete() error {
	err := n.execSQL("DELETE FROM note WHERE id=:id", n)
	if err == nil {
		n.ID = 0
	}
	return err
}

//Scan implements driver Scanner interface
func (n *Note) Scan(value interface{}) error {
	n.ID = value.(int64)
	return nil
}

//Value implements the driver Valuer interface
func (n Note) Value() (driver.Value, error) {
	return n.ID, nil
}
