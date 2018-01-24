package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-oci8"
)

func getDSN() string {
	// same as "sqlplus sys/syspwd@tnsentry as sysdba"
	return "sys/syspassword@mytnsentry?as=sysdba"
}
func main() {
	os.Setenv("NLS_LANG", "")

	db, err := sql.Open("oci8", getDSN())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	fmt.Println()
	var user string
	err = db.QueryRow("select user from dual").Scan(&user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Successful 'as sysdba' connection. Current user is: %v\n", user)
}
