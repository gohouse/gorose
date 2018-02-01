package drivers

import (
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
)

// MySQL driver
func MySQL(dbObj map[string]string) (driver string, dsn string) {
	// driver
	driver = "mysql"

	// dsn string
	dsn = fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s",
		dbObj["username"], dbObj["password"], dbObj["protocol"], dbObj["host"],
		dbObj["port"], dbObj["database"], dbObj["charset"])

	return
}
