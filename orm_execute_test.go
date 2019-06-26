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
	fmt.Println(aff, err)
	fmt.Println(db.LastSql())
}

func Test_Transaction(t *testing.T) {
	var db = initOrm()
	// 一键事务, 自动回滚和提交, 我们只需要关注业务即可
	err := db.Transaction(
		func(db IOrm) error {
			_,err := db.Where("uid",1).Update(&Data{"name":"gorose2"})
			if err!=nil {
				return err
			}
			_,err = db.Insert(&Data{"name":"gorose2"})
			if err!=nil {
				return err
			}
			return nil
		},
		func(db IOrm) error {
			_,err := db.Where("uid",1).Delete()
			if err!=nil {
				return err
			}
			_,err = db.Insert(&Data{"name":"gorose2"})
			if err!=nil {
				return err
			}
			return nil
		})
	t.Log(err)
}

func trans_func1(db IOrm) error {
	return nil
}

func trans_func2(db IOrm) error {
	return nil
}
