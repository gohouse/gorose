package gorose

import (
	"fmt"
	"testing"
)

func TestOrm_Update(t *testing.T) {
	db := DB()

	var u = Users{
		Name: "gorose2",
		Age:  19,
	}

	aff, err := db.Force().Data(u).Update()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(aff, db.LastSql())
}

func TestOrm_UpdateMap(t *testing.T) {
	db := DB()

	var u = &UsersMapSlice{{"name": "gorose2", "age": 19}}

	aff, err := db.Force().Data(u).Update()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(aff, db.LastSql())
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
