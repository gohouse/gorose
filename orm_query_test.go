package gorose

import (
	"fmt"
	"testing"
)

func TestOrm_First(t *testing.T) {
	db := initOrm()
	var u = Users{}
	res,err := db.Table(&u).Get()
	if err!=nil {
		t.Fatal(err.Error())
	}
	t.Log(res)
	db.Table(&u).Limit(2).Select()
	t.Log(u)
}

func TestOrm_Get2(t *testing.T) {
	db := initOrm()
	var u =[]Users{}
	var err error

	fmt.Println(&u)

	err = db.Table(&u).Limit(1).Select()
	fmt.Println(err, u, db.LastSql())

	//res, err := db.Table("users").Where("uid", ">", 2).Limit(2).Get()
	//fmt.Println(err, res, db.LastSql())
	//
	//fmt.Println(u)
}
