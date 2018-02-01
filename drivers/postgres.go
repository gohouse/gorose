package drivers

import (
	"fmt"
	//_ "github.com/lib/pq"
)

// Postgres driver
func Postgres(dbObj map[string]string) (driver string, dsn string) {
	// driver
	driver = "postgres"

	// dsn string
	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbObj["host"], dbObj["port"], dbObj["username"], dbObj["password"], dbObj["database"])

	return
}
