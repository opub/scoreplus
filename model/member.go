package model

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
)

//Member data model
type Member struct {
	Base
	Handle     string        `sql:" NOT NULL UNIQUE"`
	Email      string        `sql:" NOT NULL UNIQUE"`
	FirstName  string        `json:"firstname"`
	LastName   string        `json:"lastname"`
	Verified   bool          `json:"-"`
	Enabled    bool          `json:"-"`
	LastActive null.Time     `json:"-"`
	Teams      pq.Int64Array `json:"teams,omitempty"`
	Follows    pq.Int64Array `json:"follows,omitempty"`
	Followers  pq.Int64Array `json:"followers,omitempty"`
}

//Save persists object to data store
func (m *Member) Save() error {
	if m.ID == 0 {
		m.Created = nullTimeNow()
		return m.execSQL("INSERT INTO member (handle, email, firstname, lastname, verified, enabled, lastactive, teams, follows, followers, created, createdby) VALUES (:handle, :email, :firstname, :lastname, :verified, :enabled, :lastactive, :teams, :follows, :followers, :created, :createdby) RETURNING id", m)
	}
	m.Modified = nullTimeNow()
	return m.execSQL("UPDATE member SET handle=:handle, email=:email, firstname=:firstname, lastname=:lastname, verified=:verified, enabled=:enabled, lastactive=:lastactive, teams=:teams, follows=:follows, followers=:followers, modified=:modified, modifiedby=:modifiedby WHERE id=:id", m)
}

//Delete removes object from data store
func (m *Member) Delete() error {
	return m.delete("member")
}

//SelectMembers from data store where ID in slice
func SelectMembers(ids []int64) ([]Member, error) {
	rows, err := selectRows(ids, "member")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Member, 0)
	for rows.Next() {
		m := Member{}
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

//SelectAllMembers from data store
func SelectAllMembers() ([]Member, error) {
	return SelectMembers(nil)
}
