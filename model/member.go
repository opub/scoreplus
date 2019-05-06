package model

import (
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/opub/scoreplus/db"
	"github.com/opub/scoreplus/util"
	"github.com/rs/zerolog/log"
)

//Member data model
type Member struct {
	Base
	Handle        string        `sql:" NOT NULL UNIQUE"`
	Email         string        `sql:" NOT NULL UNIQUE"`
	FirstName     string        `json:"firstname"`
	LastName      string        `json:"lastname"`
	Provider      string        `json:"-"`
	ProviderID    string        `json:"-"`
	Enabled       bool          `json:"-"`
	LastActive    null.Time     `json:"-"`
	Teams         pq.Int64Array `json:"-"`
	TeamCount     int           `sql:"-" json:"teamCount"`
	Follows       pq.Int64Array `json:"-"`
	FollowCount   int           `sql:"-" json:"followCount"`
	Followers     pq.Int64Array `json:"-"`
	FollowerCount int           `sql:"-" json:"followerCount"`
}

//Save persists object to data store
func (m *Member) Save() error {
	m.setup()
	if m.ID == 0 {
		m.Created = NullTimeNow()
		return m.execSQL("INSERT INTO member (handle, email, firstname, lastname, provider, providerid, enabled, lastactive, teams, follows, followers, created, createdby) VALUES (:handle, :email, :firstname, :lastname, :provider, :providerid, :enabled, :lastactive, :teams, :follows, :followers, :created, :createdby) RETURNING id", m)
	}
	m.Modified = NullTimeNow()
	return m.execSQL("UPDATE member SET handle=:handle, email=:email, firstname=:firstname, lastname=:lastname, provider=:provider, providerid=:providerid, enabled=:enabled, lastactive=:lastactive, teams=:teams, follows=:follows, followers=:followers, modified=:modified, modifiedby=:modifiedby WHERE id=:id", m)
}

//Delete removes object from data store
func (m *Member) Delete() error {
	return m.delete("member")
}

func SearchMembers(search string) ([]Member, error) {
	db, err := db.Connect()
	if err != nil {
		return nil, err
	}

	term := strings.ToLower(strings.TrimSpace(search))
	sql := "SELECT * FROM member WHERE lower(handle) LIKE '%' || $1 || '%' OR lower(firstname) LIKE '%' || $1 || '%' OR lower(lastname) LIKE '%' || $1 || '%'"
	rows, err := db.Queryx(sql, term)
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

//GetMember returns member from data store
func GetMember(id int64) (Member, error) {
	m := Member{}
	err := get(id, &m)
	m.setup()
	return m, err
}

//GetMemberFromProvider returns a member from data store based on provider name and userID
func GetMemberFromProvider(name string, id string) (Member, error) {
	m := Member{}

	db, err := db.Connect()
	if err != nil {
		return m, err
	}

	sql := "SELECT * FROM member WHERE provider=$1 AND providerid=$2 LIMIT 1"

	log.Info().Str("provider", name).Str("id", id).Msg("member by provider")

	rows, err := db.Queryx(sql, name, id)
	if err != nil {
		return m, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&m)
		m.setup()
	}
	return m, err
}

func (m *Member) setup() {
	m.TeamCount = len(m.Teams)
	m.FollowCount = len(m.Follows)
	m.FollowerCount = len(m.Followers)

	//force a handle based on email address if missing
	if m.Email != "" && m.Handle == "" {
		at := strings.Index(m.Email, "@")
		name := m.Email[:at]
		m.Handle = fmt.Sprintf("%s_%s", name, util.RandomString(6))
	}
}
