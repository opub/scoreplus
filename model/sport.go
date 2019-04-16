package model

//Sport data model
type Sport struct {
	Base
	Name string
}

//Save persists object to data store
func (s *Sport) Save() error {
	if s.ID == 0 {
		s.Created = nullTimeNow()
		return s.execSQL("INSERT INTO sport (name, created, createdby) VALUES (:name, :created, :createdby) RETURNING id", s)
	}
	s.Modified = nullTimeNow()
	return s.execSQL("UPDATE sport SET name=:name, modified=:modified, modifiedby=:modifiedby WHERE id=:id", s)
}

//Delete removes object from data store
func (s *Sport) Delete() error {
	return s.delete("sport")
}
