package gorose

import (
	"fmt"
	"testing"
)

func initOrm() IOrm {
	return initDB().NewOrm()
}
func DB() IOrm {
	return initDB().NewOrm()
}
func TestNewOrm(t *testing.T) {
	orm := initOrm()
	orm.Hello()
}
func TestOrm_AddFields(t *testing.T) {
	orm := initOrm()
	var u = Users{}
	var fieldStmt = orm.Table(&u).Fields("a").Where("m", 55)
	a, b, c := fieldStmt.AddFields("b").Where("d", 1).BuildSql()
	fmt.Println(a, b, c)

	d, e, f := fieldStmt.AddFields("c").Where("d", 2).BuildSql()
	fmt.Println(d, e, f)
}

func TestOrm_BuildSql(t *testing.T) {
	var u = Users{
		Name: "gorose2",
		Age:  19,
	}

	//aff, err := db.Force().Data(&u)
	a, b, c := DB().Table(&u).Where("age", ">", 1).Data(&u).BuildSql("update")
	fmt.Println(a, b, c)
}
