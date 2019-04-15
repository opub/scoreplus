package model

import "database/sql/driver"

//Sport data model
type Sport struct {
	Base
	Name string
}

//Save persists object to data store
func (s *Sport) Save() error {
	if s.ID == 0 {
		s.Created = NullTimeNow()
		return s.execSQL("INSERT INTO sport (name, created, createdby) VALUES (:name, :created, :createdby) RETURNING id", s)
	}
	s.Modified = NullTimeNow()
	return s.execSQL("UPDATE sport SET name=:name, modified=:modified, modifiedby=:modifiedby WHERE id=:id", s)
}

//Delete removes object from data store
func (s *Sport) Delete() error {
	err := s.execSQL("DELETE FROM sport WHERE id=:id", s)
	if err == nil {
		s.ID = 0
	}
	return err
}

//Scan implements driver Scanner interface
func (s *Sport) Scan(value interface{}) error {
	s.ID = value.(int64)
	return nil
}

//Value implements the driver Valuer interface
func (s Sport) Value() (driver.Value, error) {
	return s.ID, nil
}
