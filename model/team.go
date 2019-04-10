package model

//Team data model
type Team struct {
	Base
	Name   string
	Sport  Sport
	Venue  Venue
	Mascot string
	Games  []int
}
