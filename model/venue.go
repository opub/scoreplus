package model

import "database/sql/driver"

//Venue data model
type Venue struct {
	Base
	Name        string
	Address     string
	Coordinates string
}

//Save persists object to data store
func (v *Venue) Save() error {
	if v.ID == 0 {
		v.Created = NullTimeNow()
		return v.execSQL("INSERT INTO venue (name, address, coordinates, created, createdby) VALUES (:name, :address, :coordinates, :created, :createdby) RETURNING id", v)
	}
	v.Modified = NullTimeNow()
	return v.execSQL("UPDATE venue SET name=:name, address=:address, coordinates=:coordinates, modified=:modified, modifiedby=:modifiedby WHERE id=:id", v)
}

//Delete removes object from data store
func (v *Venue) Delete() error {
	err := v.execSQL("DELETE FROM venue WHERE id=:id", v)
	if err == nil {
		v.ID = 0
	}
	return err
}

//Scan implements driver Scanner interface
func (v *Venue) Scan(value interface{}) error {
	v.ID = value.(int64)
	return nil
}

//Value implements the driver Valuer interface
func (v Venue) Value() (driver.Value, error) {
	return v.ID, nil
}
