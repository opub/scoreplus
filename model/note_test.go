package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNoteCRUD(t *testing.T) {
	//create
	message := random()
	m1 := Note{Message: message}
	err := m1.Save()
	if err != nil {
		t.Errorf("insert failed: %v", err)
	}
	if m1.ID == 0 {
		t.Errorf("id not set after insert")
	}
	if m1.Message != message {
		t.Errorf("data not persisted on insert")
	}

	//read
	id := m1.ID
	m2 := Note{}
	err = Get(id, &m2)
	if err != nil {
		t.Errorf("select failed: %v", err)
	}
	if !cmp.Equal(m1, m2) {
		t.Errorf("read data doesn't match:\nm1: %+v\nm2: %+v", m1, m2)
	}

	//update
	message = random()
	m1.Message = message
	err = m1.Save()
	if err != nil {
		t.Errorf("update failed: %v", err)
	}
	if m1.ID != id {
		t.Errorf("id changed during update")
	}
	if m1.Message != message {
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

func TestNoteSelect(t *testing.T) {
	n1 := testNote()
	n2 := testNote()
	n3 := testNote()
	expected := []Note{n1, n2, n3}

	results, err := SelectNotes([]int64{n1.ID, n2.ID, n3.ID})
	if err != nil {
		t.Errorf("select failed: %v", err)
	}

	if !cmp.Equal(results, expected) {
		t.Errorf("select results don't match:\nexpected: %+v\nresults: %+v", expected, results)
	}

	n1.Delete()
	n2.Delete()
	n3.Delete()
}

func testNote() Note {
	n := Note{Message: random()}
	n.Save()
	return n
}
