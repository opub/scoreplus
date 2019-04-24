package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMemberCRUD(t *testing.T) {
	//create
	handle := random()
	email := handle + "@score.plus"
	m1 := Member{Handle: handle, Email: email, FirstName: "John", LastName: "Doe"}
	err := m1.Save()
	if err != nil {
		t.Errorf("insert failed: %v", err)
	}
	if m1.ID == 0 {
		t.Errorf("id not set after insert")
	}
	if m1.Handle != handle || m1.Email != email || m1.FirstName != "John" || m1.LastName != "Doe" {
		t.Errorf("data not persisted on insert")
	}

	//read
	id := m1.ID
	m2, err := GetMember(id)
	if err != nil {
		t.Errorf("select failed: %v", err)
	}
	if !cmp.Equal(m1, m2) {
		t.Errorf("read data doesn't match:\nm1: %+v\nm2: %+v", m1, m2)
	}

	//update
	m1.FirstName = "Jane"
	err = m1.Save()
	if err != nil {
		t.Errorf("update failed: %v", err)
	}
	if m1.ID != id {
		t.Errorf("id changed during update")
	}
	if m1.FirstName != "Jane" {
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

func TestMemberSelect(t *testing.T) {
	m1 := testMember()
	defer m1.Delete()
	m2 := testMember()
	defer m2.Delete()
	m3 := testMember()
	defer m3.Delete()
	expected := []Member{m1, m2, m3}

	results, err := SelectMembers([]int64{m1.ID, m2.ID, m3.ID})
	if err != nil {
		t.Errorf("select failed: %v", err)
	}

	if !cmp.Equal(results, expected) {
		t.Errorf("select results don't match:\nexpected: %+v\nresults: %+v", expected, results)
	}

	results, err = SelectAllMembers()
	if err != nil {
		t.Errorf("select all failed: %v", err)
	}

	if len(results) < len(expected) {
		t.Errorf("select all missing results:\n%+v", results)
	}
}

func testMember() Member {
	handle := random()
	m := Member{Handle: handle, Email: handle + "@score.plus", FirstName: handle, LastName: handle}
	m.Save()
	return m
}
