package gorose

import (
	"testing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/demo/config"
)
func TestDatabase_First(test *testing.T) {
	Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db Database

	res := db.Table("users").Where("id",">",2).First()

	if  _, ok := res["id"]; ok==false {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: id=", res["id"])
	}
}
func TestDatabase_Get(test *testing.T) {
	Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db Database

	res := db.Table("users").Where("id",">",2).Limit(10).Get()

	if  _, ok := res[0]["id"]; ok==false {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: id=", res[0]["id"])
	}
}
func TestDatabase_Insert(test *testing.T) {
	Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db Database

	// insert
	res := db.Table("users").Data(map[string]interface{}{"name":"fizz5", "age":19}).Insert()

	if  res>0 {
		test.Log("PASS: insert=", res)
	} else {
		test.Error("FAIL: insert failed.")
	}

	// update
	res2 := db.Table("users").Data(map[string]interface{}{"name":"fizz6", "age":19}).Where("id",res).Update()

	if  res2>0 {
		test.Log("PASS: Update=", res2)
	} else {
		test.Error("FAIL: Update failed.")
	}

	// delete
	res3 := db.Table("users").Data(map[string]interface{}{"name":"fizz5", "age":19}).Where("id",res).Delete()

	if  res3>0 {
		test.Log("PASS: Delete=", res3)
	} else {
		test.Error("FAIL: Delete failed.")
	}
}
func TestDatabase_Count(test *testing.T) {
	Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db Database

	res := db.Table("users").Where("id",">",100).Count()

	if  res>=0 {
		test.Log("PASS: count=", res)
	} else {
		test.Error("FAIL: test failed.")
	}
}
func TestDatabase_Sum(test *testing.T) {
	Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db Database

	res := db.Table("users").Where("id",">",2).Sum("age")

	if  res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: sum=", res)
	}
}
func TestDatabase_Avg(test *testing.T) {
	Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db Database

	res := db.Table("users").Where("id",">",2).Avg("age")

	if  res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: avg=", res)
	}
}
func TestDatabase_Max(test *testing.T) {
	Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db Database

	res := db.Table("users").Where("id",">",2).Max("age")

	if  res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: max=", res)
	}
}
func TestDatabaseMin(test *testing.T) {
	Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db Database

	res := db.Table("users").Where("id",">",2).Min("age")

	if  res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: min=", res)
	}
}
func BenchmarkDatabase_First(bmtest *testing.B) {
	Connect.Open(config.Configs,"mysql_dev")
	// close DB
	defer DB.Close()
	// get the db chaining object
	var db Database
	for cnt := 0; cnt < bmtest.N; cnt++ {
		//db.Table("users").Where("id",">",2).First()	// 316623 ns
		//db.Table("users").Fields("id").First()	// 279397
		db.Table("users").Fields("id").Where("id","<",10).Group("id").Order("id asc").First()	// 319225
	}

}
