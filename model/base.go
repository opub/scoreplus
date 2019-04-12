package model

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/opub/scoreplus/db"
)

//Model data store operations
type Model interface {
	Save() error
	Delete() error
}

//Base model that provides common fields
type Base struct {
	Model
	ID         int64
	Created    null.Time `sql:" NOT NULL DEFAULT now()"`
	CreatedBy  int64     `sql:" NOT NULL DEFAULT 0"`
	Modified   null.Time
	ModifiedBy int64 `sql:" NOT NULL DEFAULT 0"`
}

//Get matching model from data store
func Get(id int64, m interface{}) error {
	db, err := db.Connect()
	if err != nil {
		return err
	}

	table := strings.ToLower(reflect.TypeOf(m).Elem().Name())
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id=%d", table, id)

	fmt.Println(sql)

	//	return db.Get(&m, sql)

	rows, err := db.Queryx(sql)
	if rows.Next() {
		err = rows.StructScan(m)
		fmt.Printf("%#v\n", m)
		return err
	}
	return nil
}

// NullTimeNow returns a time.Now() representation truncated to a microsecond to match database precision
func NullTimeNow() null.Time {
	return null.Time{Time: time.Now().Truncate(time.Microsecond), Valid: true}
}

func (b *Base) execSQL(sql string, m interface{}) error {
	db, err := db.Connect()
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(sql)
	if err != nil {
		tx.Rollback()
		return err
	}

	if strings.Contains(strings.ToUpper(sql), " RETURNING ") {
		err = stmt.Get(&b.ID, m)
	} else {
		_, err = stmt.Exec(m)
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
