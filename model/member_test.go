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
	m2 := Member{}
	err = Get(id, &m2)
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
