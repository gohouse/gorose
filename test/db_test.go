package test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/examples/config"
	"testing"
	"github.com/gohouse/gorose"
	"fmt"
)

// go test -v
// go test -test.bench=.
func TestDatabase_First(test *testing.T) {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res,err := db.Table("users").Where("id", ">", 2).First()
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, ok := res["id"]; ok == false {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: id=", res["id"])
	}
}
func TestDatabase_Get(test *testing.T) {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res,err := db.Table("users").Where("id", ">", 2).Limit(10).Get()
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, ok := res[0]["id"]; ok == false {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: id=", res[0]["id"])
	}
}
func TestDatabase_Insert(test *testing.T) {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	// insert
	res,err := db.Table("users").Data(map[string]interface{}{"name": "fizz5", "age": 19}).Insert()
	if err != nil {
		fmt.Println(err)
		return
	}

	if res > 0 {
		test.Log("PASS: insert=", res)
	} else {
		test.Error("FAIL: insert failed.")
	}

	// update
	res2,err := db.Table("users").Data(map[string]interface{}{"name": "fizz6", "age": 19}).Where("id", res).Update()
	if err != nil {
		fmt.Println(err)
		return
	}

	if res2 > 0 {
		test.Log("PASS: Update=", res2)
	} else {
		test.Error("FAIL: Update failed.")
	}

	// delete
	res3,err := db.Table("users").Data(map[string]interface{}{"name": "fizz5", "age": 19}).Where("id", res).Delete()
	if err != nil {
		fmt.Println(err)
		return
	}

	if res3 > 0 {
		test.Log("PASS: Delete=", res3)
	} else {
		test.Error("FAIL: Delete failed.")
	}
}
func TestDatabase_Count(test *testing.T) {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res,err := db.Table("users").Where("id", ">", 100).Count()
	if err != nil {
		fmt.Println(err)
		return
	}

	if res >= 0 {
		test.Log("PASS: count=", res)
	} else {
		test.Error("FAIL: test failed.")
	}
}
func TestDatabase_Sum(test *testing.T) {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res,err := db.Table("users").Where("id", ">", 2).Sum("age")
	if err != nil {
		fmt.Println(err)
		return
	}

	if res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: sum=", res)
	}
}
func TestDatabase_Avg(test *testing.T) {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res,err := db.Table("users").Where("id", ">", 2).Avg("age")
	if err != nil {
		fmt.Println(err)
		return
	}

	if res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: avg=", res)
	}
}
func TestDatabase_Max(test *testing.T) {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res,err := db.Table("users").Where("id", ">", 2).Max("age")
	if err != nil {
		fmt.Println(err)
		return
	}

	if res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: max=", res)
	}
}
func TestDatabase_Min(test *testing.T) {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	res,err := db.Table("users").Where("id", ">", 2).Min("age")
	if err != nil {
		fmt.Println(err)
		return
	}

	if res == nil {
		test.Error("FAIL: test failed.")
	} else {
		test.Log("PASS: min=", res)
	}
}
func BenchmarkDatabase_First(bmtest *testing.B) {
	db, err := gorose.Open(config.DbConfig, "mysql_dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// close DB
	defer db.Close()

	for cnt := 0; cnt < bmtest.N; cnt++ {
		//db.Table("users").Where("id",">",2).First()	// 316623 ns
		db.Table("users").Fields("id").First() // 279397
		//db.Table("users").Fields("id").Where("id","<",10).Group("id").Order("id asc").First()	// 319225
	}

}
//func BenchmarkDatabase_Get(bmtest *testing.B) {
//	db, err := gorose.Open(config.DbConfig, "mysql_dev")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	// close DB
//	defer db.Close()
//
//	for cnt := 0; cnt < bmtest.N; cnt++ {
//		db.Table("users").Fields("id").Limit(10).Get() // 279397
//	}
//}
