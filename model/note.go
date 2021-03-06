package model

import "github.com/opub/scoreplus/util"

//Note data model
type Note struct {
	Base
	Message string `json:"message,omitempty"`
}

//Save persists object to data store
func (n *Note) Save() error {
	if n.ID == 0 {
		n.Created = NullTimeNow()
		return n.execSQL("INSERT INTO note (message, created, createdby) VALUES (:message, :created, :createdby) RETURNING id", n)
	}
	n.Modified = NullTimeNow()
	return n.execSQL("UPDATE note SET message=:message, modified=:modified, modifiedby=:modifiedby WHERE id=:id", n)
}

//Delete removes object from data store
func (n *Note) Delete() error {
	return n.delete("note")
}

//LinkID gets ID as a linkable string
func (n Note) LinkID() string {
	return util.EncodeLink(n.ID, 30)
}

//SelectNotes from data store where ID in slice
func SelectNotes(ids []int64) ([]Note, error) {
	rows, err := selectRows(ids, "note")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Note, 0)
	for rows.Next() {
		n := Note{}
		err = rows.StructScan(&n)
		if err != nil {
			return nil, err
		}
		results = append(results, n)
	}
	return results, nil
}

//SelectAllNotes from data store
func SelectAllNotes() ([]Note, error) {
	return SelectNotes(nil)
}

//GetNote returns note from data store
func GetNote(id int64) (Note, error) {
	n := Note{}
	err := get(id, &n)
	return n, err
}
