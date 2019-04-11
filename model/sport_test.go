package model

import "testing"

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
		t.Errorf("name not persisted on insert")
	}

	//read
	id := m1.ID
	m2 := Sport{}
	err = Get(id, &m2)
	if err != nil {
		t.Errorf("select failed: %v", err)
	}
	// if m1 != m2 {
	if m1.ID != m2.ID || m1.Name != m2.Name {
		t.Errorf("read data doesn't match")
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
