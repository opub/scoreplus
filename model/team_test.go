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
		t.Errorf("read data doesn't match:\nm1: %+v\nm2: %+v", m1, m2)
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

func TestTeamSelect(t *testing.T) {
	t1 := testTeam()
	defer t1.Delete()
	t2 := testTeam()
	defer t2.Delete()
	t3 := testTeam()
	defer t3.Delete()
	expected := []Team{t1, t2, t3}

	results, err := SelectTeams([]int64{t1.ID, t2.ID, t3.ID})
	if err != nil {
		t.Errorf("select failed: %v", err)
	}

	if !cmp.Equal(results, expected) {
		t.Errorf("select results don't match:\nexpected: %+v\nresults: %+v", expected, results)
	}

	results, err = SelectAllTeams()
	if err != nil {
		t.Errorf("select all failed: %v", err)
	}

	if len(results) < len(expected) {
		t.Errorf("select all missing results:\n%+v", results)
	}
}

func testTeam() Team {
	t := Team{Name: random(), Mascot: random()}
	t.Save()
	return t
}
