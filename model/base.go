package model

import (
	"database/sql/driver"
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

//Delete removes object from data store
func (b *Base) delete(table string) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE id=:id", table)
	fmt.Println(sql)
	err := b.execSQL(sql, b)
	if err == nil {
		b.ID = 0
	}
	return err
}

//Get matching model from data store
func Get(id int64, model interface{}) error {
	db, err := db.Connect()
	if err != nil {
		return err
	}

	table := strings.ToLower(reflect.TypeOf(model).Elem().Name())
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id=%d", table, id)

	rows, err := db.Queryx(sql)
	defer rows.Close()

	if rows.Next() {
		return rows.StructScan(model)
	}
	return nil
}

func createSlice(t reflect.Type) interface{} {
	var sliceType reflect.Type
	sliceType = reflect.SliceOf(t)
	return reflect.Zero(sliceType).Interface()
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

//Scan implements driver Scanner interface
func (b *Base) Scan(value interface{}) error {
	b.ID = value.(int64)
	return nil
}

//Value implements the driver Valuer interface
func (b Base) Value() (driver.Value, error) {
	return b.ID, nil
}

func table(t reflect.Type) string {
	return strings.ToLower(t.Name())
}

func nullTimeNow() null.Time {
	return null.Time{Time: time.Now().Truncate(time.Microsecond), Valid: true}
}
