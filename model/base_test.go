package model

import (
	"github.com/guregu/null"
	"github.com/opub/scoreplus/util"
)

//creates new models to get a valid ID but then clear other fields for simple equality tests

func getSport() Sport {
	s := Sport{Name: random()}
	s.Save()
	s.Name = ""
	s.Created = null.Time{}
	s.CreatedBy = 0
	return s
}

func getTeam() Team {
	t := Team{Name: random(), Mascot: random()}
	t.Save()
	t.Name = ""
	t.Mascot = ""
	t.Created = null.Time{}
	t.CreatedBy = 0
	return t
}

func random() string {
	value, _ := util.RandomID()
	return value
}
