package main

import (
	"database/sql"
	"fmt"
	_ "github.com/gohouse/gorose/driver/mysql"
	"github.com/gohouse/gorose/utils"
	"reflect"
	"time"
)

type users4 struct {
	Id int64 `orm:"id"`
	Name string `orm:"name"`
	Age int `orm:"age"`
	Job *string `orm:"job"`
	CreatedAt *time.Time `orm:"created_at"`
}


func main() {
	db, err := sql.Open("mysql",
		"gcore:gcore@tcp(192.168.200.248:3306)/test?parseTime=true")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows,err := db.Query("select id,name,age,job,created_at from users where id=55")

	defer rows.Close()

	var u users4

	for rows.Next() {
		err := rows.Scan(&u.Id,&u.Name, &u.Age, &u.Job, &u.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(u)
		fmt.Println((u.Job)==nil)
		//fmt.Println(string(u.Job))
		fmt.Println(reflect.TypeOf(u.CreatedAt))
		fmt.Println(u.CreatedAt.Format(utils.DATETIME_FORMAT))
	}

}
