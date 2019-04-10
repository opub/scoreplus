package model

import "time"

//Model data store operations
type Model interface {
	Save() error
	Insert() error
	Update() error
	Delete() error
}

//Base model that provides common fields
type Base struct {
	Model
	ID         int
	Created    time.Time `sql:" NOT NULL DEFAULT now()"`
	CreatedBy  int       `sql:" NOT NULL"`
	Modified   time.Time
	ModifiedBy int
}

//Save persists object to data store
func (b *Base) Save() error {
	if b.ID == 0 {
		return b.Insert()
	}
	return b.Update()
}
