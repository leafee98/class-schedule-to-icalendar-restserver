package db

import (
	"github.com/jmoiron/sqlx"
	// for the side effect of driver register
	_ "github.com/go-sql-driver/mysql"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/config"
)

// DB is database object to do query
var DB *sqlx.DB

// Init initialize the db, make it usable
func Init() error {
	var err error
	// auto parse database type datetime as []uint8 in go to time.Time
	DB, err = sqlx.Connect("mysql", config.DatabaseSource+"?parseTime=true")
	return err
}
