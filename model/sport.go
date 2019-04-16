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

//SelectSports from data store where ID in slice
func SelectSports(ids []int64) ([]Sport, error) {
	rows, err := selectRows(ids, "sport")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Sport, 0)
	for rows.Next() {
		s := Sport{}
		err = rows.StructScan(&s)
		if err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, nil
}

//SelectAllSports from data store
func SelectAllSports() ([]Sport, error) {
	return SelectSports(nil)
}
