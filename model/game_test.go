package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var sport = getSport()
var team1 = getTeam()
var team2 = getTeam()

func TestGameCRUD(t *testing.T) {
	defer cleanup()

	//create
	m1 := Game{Sport: sport, HomeTeam: team1, AwayTeam: team2, HomeScore: 10, AwayScore: 1}
	err := m1.Save()
	if err != nil {
		t.Errorf("insert failed: %v", err)
	}
	if m1.ID == 0 {
		t.Errorf("id not set after insert")
	}
	if m1.HomeScore != 10 || m1.AwayScore != 1 {
		t.Errorf("data not persisted on insert")
	}

	//read
	id := m1.ID
	m2 := Game{}
	err = Get(id, &m2)
	if err != nil {
		t.Errorf("select failed: %v", err)
	}
	if !cmp.Equal(m1, m2) {
		t.Errorf("read data doesn't match:\nm1: %+v\nm2: %+v", m1, m2)
	}

	//update
	m1.HomeScore = 11
	m1.AwayScore = 2
	err = m1.Save()
	if err != nil {
		t.Errorf("update failed: %v", err)
	}
	if m1.ID != id {
		t.Errorf("id changed during update")
	}
	if m1.HomeScore != 11 || m1.AwayScore != 2 {
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

func cleanup() {
	sport.Delete()
	team1.Delete()
	team2.Delete()
}
