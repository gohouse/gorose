package drivers

import (
//_ "github.com/mattn/go-sqlite3"
)

// Sqlite3 driver
func Sqlite3(dbObj map[string]string) (driver string, dsn string) {
	// driver
	driver = "sqlite3"

	// dsn string
	dsn = dbObj["database"]

	return
}
