package db

import (
	"fmt"

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
	DB, err = sqlx.Connect("mysql", dsn(config.DatabaseUsername, config.DatabasePassword,
		config.DatabaseHost, config.DatabaseName))
	return err
}

func dsn(username, password, host, name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, name)
}
