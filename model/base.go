package model

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"

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
	Model      `json:"-"`
	ID         int64     `json:"id,omitempty"`
	Created    null.Time `sql:" NOT NULL DEFAULT now()" json:"created,omitempty"`
	CreatedBy  int64     `sql:" NOT NULL DEFAULT 0" json:"createdby,omitempty"`
	Modified   null.Time `json:"modified,omitempty"`
	ModifiedBy int64     `sql:" NOT NULL DEFAULT 0" json:"modifiedby,omitempty"`
}

//Delete removes object from data store
func (b *Base) delete(table string) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE id=:id", table)
	err := b.execSQL(sql, b)
	if err == nil {
		b.ID = 0
	}
	return err
}

//get matching model from data store
func get(id int64, model interface{}) error {
	db, err := db.Connect()
	if err != nil {
		return err
	}

	table := strings.ToLower(reflect.TypeOf(model).Elem().Name())
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id=$1 LIMIT 1", table)

	log.Info().Str("table", table).Int64("id", id).Msg("model.get")

	rows, err := db.Queryx(sql, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.StructScan(model)
	}
	return nil
}

//get models from data store if ids is empty then all rows returned
func selectRows(ids []int64, table string) (*sqlx.Rows, error) {
	db, err := db.Connect()
	if err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		sql := fmt.Sprintf("SELECT * FROM %s WHERE id = ANY($1)", table)
		return db.Queryx(sql, pq.Array(ids))
	}
	sql := fmt.Sprintf("SELECT * FROM %s", table)
	return db.Queryx(sql)
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

//NullTimeNow returns time.Now() as a nullable database value
func NullTimeNow() null.Time {
	return WrapTime(time.Now())
}

//WrapTime wraps a time.Time value as a null.Time for database usage
func WrapTime(t time.Time) null.Time {
	return null.Time{Time: t.Truncate(time.Microsecond), Valid: true}
}
