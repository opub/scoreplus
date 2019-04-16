package model

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
)

//Member data model
type Member struct {
	Base
	Handle     string `sql:" NOT NULL UNIQUE"`
	Email      string `sql:" NOT NULL UNIQUE"`
	FirstName  string
	LastName   string
	Verified   bool
	Enabled    bool
	LastActive null.Time
	Teams      pq.Int64Array
	Follows    pq.Int64Array
	Followers  pq.Int64Array
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
