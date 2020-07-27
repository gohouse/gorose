package gorose

import (
	"fmt"
	"testing"
)

func TestOrm_Update(t *testing.T) {
	db := DB()

	var u = []Users{{
		Name: "gorose2",
		Age:  19,
	}}

	aff, err := db.Force().Update(&u)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(aff, db.LastSql())
}

func TestOrm_Update2(t *testing.T) {
	db := DB()

	//var u = []Users{{
	//	Name: "gorose2",
	//	Age:  11,
	//}}

	aff, err := db.Table("users").Where("uid", 1).Update()
	if err != nil {
		//t.Error(err.Error())
		t.Log(err.Error())
		return
	}
	t.Log(aff, db.LastSql())
}

func TestOrm_UpdateMap(t *testing.T) {
	db := DB()

	//var u = []UsersMap{{"name": "gorose2", "age": 19}}
	var u = UsersMap{"name": "gorose2", "age": 19}

	aff, err := db.Force().Update(&u)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(aff, db.LastSql())
}

func TestTrans(t *testing.T) {
	var db = DB()
	var db2 = DB()
	var res Users
	db.Begin()
	db2.Table(&res).Select()
	t.Log(res)
	db.Commit()
	t.Log(res)
}

func Test_Transaction(t *testing.T) {
	var db = DB()
	// 一键事务, 自动回滚和提交, 我们只需要关注业务即可
	err := db.Transaction(
		func(db IOrm) error {
			//db.Table("users").Limit(2).SharedLock().Get()
			//fmt.Println(db.LastSql())
			_, err := db.Table("users").Where("uid", 2).Update(Data{"name": "gorose2"})
			fmt.Println(db.LastSql())
			if err != nil {
				return err
			}
			_, err = db.Insert(&UsersMap{"name": "gorose2", "age": 0})
			fmt.Println(db.LastSql())
			if err != nil {
				return err
			}
			return nil
		},
		func(db IOrm) error {
			_, err := db.Table("users").Where("uid", 3).Update(Data{"name": "gorose3"})
			fmt.Println(db.LastSql())
			if err != nil {
				return err
			}
			_, err = db.Insert(&UsersMap{"name": "gorose2", "age": 0})
			fmt.Println(db.LastSql())
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("事务测试通过")
}
