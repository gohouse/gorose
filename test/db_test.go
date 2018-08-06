package test

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	"github.com/gohouse/gorose/examples/config"
	"testing"
)

// go test -v
// go test -test.bench=.

var connection gorose.Connection
var err error

func init() {
	connection, err = gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		var test *testing.T
		test.Error("FAIL: test failed.")
		return
	}
	// close DB
	//defer db.Close()
}
func TestDatabase_Query(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Query("select * from users limit ?", 1)
	if err != nil {
		test.Error("FAIL: test failed.")
		return
	}
	test.Log(fmt.Sprintf("PASS: id=?",res))
}
func TestDatabase_Execute(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Execute("insert into users set age=?", 18)
	if err != nil {
		test.Error("FAIL: test failed.", err)
		return
	}

	if res > 0 {
		test.Log("PASS: id=", res)
	} else {
		test.Error("FAIL: test failed.", db.LastSql, res)
	}
}
func TestDatabase_Value(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Table("users").Value("name")
	if err != nil {
		test.Error("FAIL: test failed.")
		return
	}
	test.Log("PASS: name =", res)
}
func TestDatabase_First(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Table("users").First()
	if err != nil {
		test.Error("FAIL: test failed.")
		return
	}

	if _, ok := res["id"]; ok == false {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: id=", res["id"])
	}
}
func TestDatabase_Get(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Table("users").Fields("id,sum(age) as sum").Where("id", ">", 0).OrWhere("age", ">", 0).
		Group("id").Having("sum>1").Order("id asc").Limit(1).Offset(0).Get()
	if err != nil {
		test.Errorf("FAIL: test failed. %s", err)
		return
	}

	//if _, ok := res[0]["id"]; ok == false {
	//	test.Error("FAIL: test failed.")
	//} else {
	//	test.Log("PASS: id=", res[0]["id"])
	//}
	test.Log(fmt.Sprintf("PASS: id=?",res))
}
func TestDatabase_Join(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Table("users a").
		Join("area b", "a.id", "=", "b.uid").Get()
	if err != nil {
		test.Error("FAIL: test failed.")
		return
	}

	//if _, ok := res[0]["id"]; ok == false {
	//	test.Error("FAIL: test failed.")
	//} else {
	//	test.Log("PASS: id=", res[0]["id"])
	//}
	test.Log(fmt.Sprintf("PASS: id=?",res))
}
func TestDatabase_InsertUpdateDelete(test *testing.T) {
	db := connection.GetInstance()
	db.Begin()
	// insert
	res, err := db.Table("users").Data(map[string]interface{}{"name": "fizz5", "age": 19}).Insert()
	if err != nil {
		db.Rollback()
		test.Error("FAIL: test failed.")
		return
	}

	if res > 0 {
		test.Log("PASS: insert=", res)
	} else {
		db.Rollback()
		test.Error("FAIL: insert failed.")
	}

	// update
	res2, err := db.Table("users").Data(map[string]interface{}{"name": "fizz6", "age": 19}).Where("id", db.LastInsertId).Update()
	if err != nil {
		db.Rollback()
		test.Error("FAIL: test failed.",db.LastSql)
		return
	}

	if res2 > 0 {
		test.Log("PASS: Update=", res2)
	} else {
		db.Rollback()
		test.Error("FAIL: test failed.",db.LastSql)
	}

	// delete
	res3, err := db.Table("users").Data(map[string]interface{}{"name": "fizz5", "age": 19}).Where("id", db.LastInsertId).Delete()
	if err != nil {
		db.Rollback()
		test.Error("FAIL: test failed.")
		return
	}

	if res3 > 0 {
		test.Log("PASS: Delete=", res3)
	} else {
		db.Rollback()
		test.Error("FAIL: Delete failed.")
	}
	db.Commit()
}
func TestDatabase_Count(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Table("users").Where("id", ">", 100).Count()
	if err != nil {
		test.Error("FAIL: test failed.")
		return
	}

	if res >= 0 {
		test.Log("PASS: count=", res)
	} else {
		test.Error("FAIL: test failed.")
	}
}
func TestDatabase_Sum(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Table("users").Where("id", ">", 2).Sum("age")
	if err != nil {
		test.Error("FAIL: test failed.")
		return
	}

	if res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: sum=", res)
	}
}
func TestDatabase_Avg(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Table("users").Where("id", ">", 2).Avg("age")
	if err != nil {
		test.Error("FAIL: test failed.")
		return
	}

	if res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: avg=", res)
	}
}
func TestDatabase_Max(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Table("users").Where("id", ">", 2).Max("age")
	if err != nil {
		test.Error("FAIL: test failed.")
		return
	}

	if res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: max=", res)
	}
}
func TestDatabase_Min(test *testing.T) {
	db := connection.GetInstance()
	res, err := db.Table("users").Where("id", ">", 2).Min("age")
	if err != nil {
		test.Error("FAIL: test failed.")
		return
	}

	if res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: min=", res)
	}
}
func BenchmarkDatabase_First(bmtest *testing.B) {
	db := connection.GetInstance()
	for cnt := 0; cnt < bmtest.N; cnt++ {
		//db.Table("users").Where("id",">",2).First()	// 316623 ns
		db.Table("users").Fields("id").First() // 279397
		//db.Table("users").Fields("id").Where("id","<",10).Group("id").Order("id asc").First()	// 319225
	}
}

//func BenchmarkDatabase_Get(bmtest *testing.B) {
//	db, err := gorose.Open(config.DbConfig, "mysql_dev")
//	if err != nil {
//		test.Error("FAIL: test failed.")
//		return
//	}
//	// close DB
//	defer db.Close()
//
//	for cnt := 0; cnt < bmtest.N; cnt++ {
//		db.Table("users").Fields("id").Limit(10).Get() // 279397
//	}
//}
