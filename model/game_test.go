package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/guregu/null"
)

func TestGameCRUD(t *testing.T) {
	s := testSimpleSport()
	defer s.Delete()
	t1 := testSimpleTeam()
	defer t1.Delete()
	t2 := testSimpleTeam()
	defer t2.Delete()

	//create
	m1 := Game{Sport: s, HomeTeam: t1, AwayTeam: t2, HomeScore: 10, AwayScore: 1}
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
	m2, err := GetGame(id)
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

func TestGameSelect(t *testing.T) {
	g1 := testGame()
	defer g1.Delete()
	g2 := testGame()
	defer g2.Delete()
	g3 := testGame()
	defer g3.Delete()
	expected := []Game{g1, g2, g3}

	results, err := SelectGames([]int64{g1.ID, g2.ID, g3.ID})
	if err != nil {
		t.Errorf("select failed: %v", err)
	}

	if !cmp.Equal(results, expected) {
		t.Errorf("select results don't match:\nexpected: %+v\nresults: %+v", expected, results)
	}

	results, err = SelectAllGames()
	if err != nil {
		t.Errorf("select all failed: %v", err)
	}

	if len(results) < len(expected) {
		t.Errorf("select all missing results:\n%+v", results)
	}
}

func testGame() Game {
	g := Game{HomeScore: 101, AwayScore: 99, Final: true, Start: NullTimeNow()}
	g.Save()
	return g
}

func testSimpleSport() Sport {
	s := testSport()
	s.Name = ""
	s.Created = null.Time{}
	s.CreatedBy = 0
	return s
}

func testSimpleTeam() Team {
	t := testTeam()
	t.Name = ""
	t.Mascot = ""
	t.Created = null.Time{}
	t.CreatedBy = 0
	return t
}
