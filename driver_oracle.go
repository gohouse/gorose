package gorose

import (
	"database/sql"
	//_ "github.com/mattn/go-oci8"
	"fmt"
)

func Oracle() {
	dbObj := CurrentConfig
	var err error

	dsn := fmt.Sprintf("%s/%s@%s", dbObj["username"], dbObj["password"], dbObj["database"])
	DB, err = sql.Open("oci8", dsn)

	CheckErr(err)
}
