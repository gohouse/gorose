package drivers

import (
	//_ "github.com/mattn/go-oci8"
	"fmt"
)

// Oracle driver
func Oracle(dbObj map[string]string) (driver string, dsn string) {
	// driver
	driver = "oci8"

	// dsn string
	dsn = fmt.Sprintf("%s/%s@%s",
		dbObj["username"], dbObj["password"], dbObj["database"])

	return
}
