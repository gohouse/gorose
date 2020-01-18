package gorose

import (
	"testing"
)

func DB() IOrm {
	return initDB().NewOrm()
}
func TestNewOrm(t *testing.T) {
	orm := DB()
	orm.Close()
}
func TestOrm_AddFields(t *testing.T) {
	orm := DB()
	//var u = Users{}
	var fieldStmt = orm.Table("users").Fields("a").Where("m", 55)
	a, b, err := fieldStmt.AddFields("b").Where("d", 1).BuildSql()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(a, b)

	fieldStmt.Reset()
	d, e, err := fieldStmt.Fields("a").AddFields("c").Where("d", 2).BuildSql()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(d, e)
}

func TestOrm_BuildSql(t *testing.T) {
	var u = Users{
		Name: "gorose2",
		Age:  19,
	}

	//aff, err := db.Force().Data(&u)
	a, b, err := DB().Table(&u).Where("age", ">", 1).Data(&u).BuildSql("update")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(a, b)
}

func TestOrm_BuildSql_where(t *testing.T) {
	var u = Users{
		Name: "gorose2",
		Age:  19,
	}

	var db = DB()
	a, b, err := db.Table(&u).Where("age", ">", 1).Where(func() {
		db.Where("name", "like", "%fizz%").OrWhere(func() {
			db.Where("age", ">", 10).Where("uid", ">", 2)
		})
	}).Limit(2).Offset(2).BuildSql()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(a, b)
}
