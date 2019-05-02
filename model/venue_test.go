package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opub/scoreplus/util"
)

func TestVenueCRUD(t *testing.T) {
	//create
	m1 := Venue{Name: "Old Yankee Stadium", Address: "East 161st Street & River Avenue"}
	err := m1.Save()
	if err != nil {
		t.Errorf("insert failed: %v", err)
	}
	if m1.ID == 0 {
		t.Errorf("id not set after insert")
	}
	if m1.Name != "Old Yankee Stadium" || m1.Address != "East 161st Street & River Avenue" {
		t.Errorf("data not persisted on insert")
	}

	//read
	id := m1.ID
	m2, err := GetVenue(id)
	if err != nil {
		t.Errorf("select failed: %v", err)
	}
	if !cmp.Equal(m1, m2) {
		t.Errorf("read data doesn't match:\nm1: %+v\nm2: %+v", m1, m2)
	}

	//update
	m1.Name = "New Yankee Stadium"
	m1.Address = "1 E 161 St, The Bronx, NY 10451"
	m1.Coordinates = "40.829167, -73.926389"
	err = m1.Save()
	if err != nil {
		t.Errorf("update failed: %v", err)
	}
	if m1.ID != id {
		t.Errorf("id changed during update")
	}
	if m1.Name != "New Yankee Stadium" || m1.Address != "1 E 161 St, The Bronx, NY 10451" || m1.Coordinates != "40.829167, -73.926389" {
		t.Errorf("name not persisted on update")
	}

	//delete
	err = m1.Delete()
	if err != nil {
		t.Errorf("delete failed: %v", err)
	}
	if m1.ID != 0 {
		t.Errorf("id not cleared during delete")
	}
}

func TestVenueSelect(t *testing.T) {
	v1 := testVenue()
	defer v1.Delete()
	v2 := testVenue()
	defer v2.Delete()
	v3 := testVenue()
	defer v3.Delete()
	expected := []Venue{v1, v2, v3}

	results, err := SelectVenues([]int64{v1.ID, v2.ID, v3.ID})
	if err != nil {
		t.Errorf("select failed: %v", err)
	}

	if !cmp.Equal(results, expected) {
		t.Errorf("select results don't match:\nexpected: %+v\nresults: %+v", expected, results)
	}

	results, err = SelectAllVenues()
	if err != nil {
		t.Errorf("select all failed: %v", err)
	}

	if len(results) < len(expected) {
		t.Errorf("select all missing results:\n%+v", results)
	}
}

func testVenue() Venue {
	v := Venue{Name: util.RandomString(10), Address: util.RandomString(10), Coordinates: util.RandomString(10)}
	v.Save()
	return v
}
