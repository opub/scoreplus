package model

import "time"

//Team data model
type Team struct {
	ID       string
	Name     string `sql:",unique"`
	Sport    Sport
	Location string
	Mascot   string
	Games    []*Game
	Created  time.Time `sql:"default:now()"`
	Modified time.Time
}
