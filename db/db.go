package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/opub/scoreplus/util"

	//loads postgresql driver
	_ "github.com/lib/pq"
)

//Connect to database
func Connect() (*sqlx.DB, error) {

	var db *sqlx.DB

	config := util.GetConfig()

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.Username, config.DB.Password, config.DB.Name)

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return db, err
	}

	return db, nil
}
