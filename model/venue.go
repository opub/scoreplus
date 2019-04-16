package model

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
		v.Created = nullTimeNow()
		return v.execSQL("INSERT INTO venue (name, address, coordinates, created, createdby) VALUES (:name, :address, :coordinates, :created, :createdby) RETURNING id", v)
	}
	v.Modified = nullTimeNow()
	return v.execSQL("UPDATE venue SET name=:name, address=:address, coordinates=:coordinates, modified=:modified, modifiedby=:modifiedby WHERE id=:id", v)
}

//Delete removes object from data store
func (v *Venue) Delete() error {
	return v.delete("venue")
}

//SelectVenues from data store where ID in slice
func SelectVenues(ids []int64) ([]Venue, error) {
	rows, err := selectRows(ids, "venue")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Venue, 0)
	for rows.Next() {
		v := Venue{}
		err = rows.StructScan(&v)
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}
	return results, nil
}

//SelectAllVenues from data store
func SelectAllVenues() ([]Venue, error) {
	return SelectVenues(nil)
}
