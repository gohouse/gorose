package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/mattn/go-oci8"
)

type ID string

func (id ID) Scan(src interface{}) error {
	fmt.Println(src)
	return nil
}

func getDSN() string {
	var dsn string
	if len(os.Args) > 1 {
		dsn = os.Args[1]
		if dsn != "" {
			return dsn
		}
	}
	dsn = os.Getenv("GO_OCI8_CONNECT_STRING")
	if dsn != "" {
		return dsn
	}
	fmt.Fprintln(os.Stderr, `Please specifiy connection parameter in GO_OCI8_CONNECT_STRING environment variable,
or as the first argument! (The format is user/name@host:port/sid)`)
	return "scott/tiger@XE"
}

func main() {
	os.Setenv("NLS_LANG", "")

	db, err := sql.Open("oci8", getDSN())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	db.Exec("drop table lastinsertid_example")

	_, err = db.Exec("create table lastinsertid_example(id varchar2(256) not null primary key, data varchar2(256))")
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := db.Exec("insert into lastinsertid_example(id, data) values(:1, :2)", "001", "こんにちわ世界")
	if err != nil {
		fmt.Println(err)
		return
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return
	}
	rowID := oci8.GetLastInsertId(lastInsertId)
	var id string
	err = db.QueryRow("select id from lastinsertid_example where rowid = :1", rowID).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(id)
}
