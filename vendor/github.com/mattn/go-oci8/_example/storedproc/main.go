package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-oci8"
)

func main() {
	db, err := sql.Open("oci8", "scott/tiger@XE")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
    CREATE OR REPLACE FUNCTION MY_SUM
    (
      P_NUM1 IN NUMBER,
      P_NUM2 IN NUMBER
    )
    RETURN NUMBER
    IS
      R_NUM NUMBER(2) DEFAULT 0;
    BEGIN
      FOR i IN 1..P_NUM2
      LOOP
        R_NUM := R_NUM + P_NUM1;
      END LOOP;
      RETURN R_NUM;
    END;
    `)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("select SCOTT.MY_SUM(5,6) from dual")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			log.Fatal(err)
		}
		println(i)
	}
}
