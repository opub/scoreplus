package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTeamCRUD(t *testing.T) {
	//create
	m1 := Team{Name: "Bad News", Mascot: "Bears", Games: []int64{0, 1, 2}}
	err := m1.Save()
	if err != nil {
		t.Errorf("insert failed: %v", err)
	}
	if m1.ID == 0 {
		t.Errorf("id not set after insert")
	}
	if m1.Name != "Bad News" || m1.Mascot != "Bears" || len(m1.Games) != 3 {
		t.Errorf("data not persisted on insert")
	}

	//read
	id := m1.ID
	m2 := Team{}
	err = Get(id, &m2)
	if err != nil {
		t.Errorf("select failed: %v", err)
	}
	if !cmp.Equal(m1, m2) {
		t.Errorf("read data doesn't match:\nm1: %v\nm2: %v", m1, m2)
	}

	//update
	m1.Name = "Good News"
	err = m1.Save()
	if err != nil {
		t.Errorf("update failed: %v", err)
	}
	if m1.ID != id {
		t.Errorf("id changed during update")
	}
	if m1.Name != "Good News" {
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
