package model

import "testing"

func TestSportInsert(t *testing.T) {
	s := Sport{Name: "Cornhole"}

	err := s.Insert()
	if err != nil {
		t.Errorf("sport insert failed: %v", err)
	}
}
