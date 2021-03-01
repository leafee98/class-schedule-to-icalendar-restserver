package db

import (
	"database/sql"

	// for the side effect of driver register
	_ "github.com/go-sql-driver/mysql"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/config"
)

// DB is database object to do query
var DB *sql.DB

// Init initialize the db, make it usable
func Init() error {
	var err error
	DB, err = sql.Open("mysql", config.DatabaseSource)
	if err != nil {
		return err
	}

	err = DB.Ping()
	return err
}
