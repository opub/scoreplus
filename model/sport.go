package model

import "github.com/opub/scoreplus/db"

//Sport data model
type Sport struct {
	Base
	Name string
}

//Insert new record
func (m *Sport) Insert() error {
	db, err := db.Connect()
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO sport (name, createdby) VALUES ($1, 1)", m.Name)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//Update existing record
func (m *Sport) Update() error {
	return nil
}
