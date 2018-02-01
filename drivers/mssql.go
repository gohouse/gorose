package drivers

import (
	"fmt"
	//_ "github.com/denisenkom/go-mssqldb"
)

// MsSQL (sqlserver) driver
func MsSQL(dbObj map[string]string) (driver string, dsn string) {
	// driver
	driver = "mssql"

	// dsn string
	dsn = fmt.Sprintf("server=%s;port=%s;database=%s;user id=%s;password=%s",
		dbObj["host"], dbObj["port"], dbObj["database"], dbObj["username"], dbObj["password"])

	return
}
