package model

import (
	"strings"

	"github.com/opub/scoreplus/db"
	"github.com/opub/scoreplus/util"
)

//Venue data model
type Venue struct {
	Base
	Name        string `json:"name,omitempty"`
	Address     string `json:"address,omitempty"`
	Coordinates string `json:"coordinates,omitempty"`
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
	return v.delete("venue")
}

//LinkID gets ID as a linkable string
func (v Venue) LinkID() string {
	return util.EncodeLink(v.ID, 50)
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

//GetVenue returns venue from data store
func GetVenue(id int64) (Venue, error) {
	v := Venue{}
	err := get(id, &v)
	return v, err
}

//SearchVenues finds venues that match a search string
func SearchVenues(search string) ([]Venue, error) {
	db, err := db.Connect()
	if err != nil {
		return nil, err
	}

	term := strings.ToLower(strings.TrimSpace(search))
	sql := "SELECT * FROM venue WHERE lower(name) LIKE '%' || $1 || '%' OR lower(address) LIKE '%' || $1 || '%'"
	rows, err := db.Queryx(sql, term)
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
