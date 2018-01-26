package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-oci8"
)

func main() {
	db, err := sql.Open("oci8", getDSN())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err = testSelect(db); err != nil {
		fmt.Println(err)
		return
	}

	if err = testI18n(db); err != nil {
		fmt.Println(err)
		return
	}
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

func testSelect(db *sql.DB) error {
	rows, err := db.Query("select 3.14, 'foo' from dual")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var f1 float64
		var f2 string
		rows.Scan(&f1, &f2)
		println(f1, f2) // 3.14 foo
	}
	return nil
}

const tbl = "tst_oci8_i18n"
const tst_strings = "'Habitación doble', '雙人房', 'двухместный номер'"

func testI18n(db *sql.DB) error {
	db, err := sql.Open("oci8", getDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	_, _ = db.Exec("DROP TABLE " + tbl)
	defer db.Exec("DROP TABLE " + tbl)
	if _, err = db.Exec("CREATE TABLE " + tbl + " (name_spainish VARCHAR2(100), name_chinesses VARCHAR2(100), name_russian VARCHAR2(100))"); err != nil {
		return err
	}
	if _, err = db.Exec("INSERT INTO " + tbl +
		" (name_spainish, name_chinesses, name_russian) " +
		" VALUES (" + tst_strings + ")"); err != nil {
		return err
	}

	rows, err := db.Query("select name_spainish, name_chinesses, name_russian from " + tbl)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var nameSpainish string
		var nameChinesses string
		var nameRussian string
		if err = rows.Scan(&nameSpainish, &nameChinesses, &nameRussian); err != nil {
			return err
		}
		got := fmt.Sprintf("'%s', '%s', '%s'", nameSpainish, nameChinesses, nameRussian)
		fmt.Println(got)
		if got != tst_strings {
			return fmt.Errorf("ERROR: string mismatch: got %q, awaited %q\n", got, tst_strings)
		}
	}
	return rows.Err()
}
