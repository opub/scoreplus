package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSportCRUD(t *testing.T) {
	//create
	m1 := Sport{Name: "Football"}
	err := m1.Save()
	if err != nil {
		t.Errorf("insert failed: %v", err)
	}
	if m1.ID == 0 {
		t.Errorf("id not set after insert")
	}
	if m1.Name != "Football" {
		t.Errorf("data not persisted on insert")
	}

	//read
	id := m1.ID
	m2 := Sport{}
	err = Get(id, &m2)
	if err != nil {
		t.Errorf("select failed: %v", err)
	}
	if !cmp.Equal(m1, m2) {
		t.Errorf("read data doesn't match:\nm1: %+v\nm2: %+v", m1, m2)
	}

	//update
	m1.Name = "Soccer"
	err = m1.Save()
	if err != nil {
		t.Errorf("update failed: %v", err)
	}
	if m1.ID != id {
		t.Errorf("id changed during update")
	}
	if m1.Name != "Soccer" {
		t.Errorf("data not persisted on update")
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

func TestSportSelect(t *testing.T) {
	s1 := testSport()
	defer s1.Delete()
	s2 := testSport()
	defer s2.Delete()
	s3 := testSport()
	defer s3.Delete()
	expected := []Sport{s1, s2, s3}

	results, err := SelectSports([]int64{s1.ID, s2.ID, s3.ID})
	if err != nil {
		t.Errorf("select failed: %v", err)
	}

	if !cmp.Equal(results, expected) {
		t.Errorf("select results don't match:\nexpected: %+v\nresults: %+v", expected, results)
	}

	results, err = SelectAllSports()
	if err != nil {
		t.Errorf("select all failed: %v", err)
	}

	if len(results) < len(expected) {
		t.Errorf("select all missing results:\n%+v", results)
	}
}

func testSport() Sport {
	s := Sport{Name: random()}
	s.Save()
	return s
}
