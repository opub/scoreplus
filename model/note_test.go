package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opub/scoreplus/util"
)

func TestNoteCRUD(t *testing.T) {
	//create
	message := util.RandomString(10)
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
	m2, err := GetNote(id)
	if err != nil {
		t.Errorf("select failed: %v", err)
	}
	if !cmp.Equal(m1, m2) {
		t.Errorf("read data doesn't match:\nm1: %+v\nm2: %+v", m1, m2)
	}

	//update
	message = util.RandomString(10)
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
	defer n1.Delete()
	n2 := testNote()
	defer n2.Delete()
	n3 := testNote()
	defer n3.Delete()
	expected := []Note{n1, n2, n3}

	results, err := SelectNotes([]int64{n1.ID, n2.ID, n3.ID})
	if err != nil {
		t.Errorf("select failed: %v", err)
	}

	if !cmp.Equal(results, expected) {
		t.Errorf("select results don't match:\nexpected: %+v\nresults: %+v", expected, results)
	}

	results, err = SelectAllNotes()
	if err != nil {
		t.Errorf("select all failed: %v", err)
	}

	if len(results) < len(expected) {
		t.Errorf("select all missing results:\n%+v", results)
	}
}

func testNote() Note {
	n := Note{Message: util.RandomString(10)}
	n.Save()
	return n
}
