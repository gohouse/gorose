package gorose

import (
	"errors"
	"testing"
)

type Users struct {
	Uid  int    `orm:"id"`
	Name string `orm:"name"`
	Age  int    `orm:"age"`
}

func (u *Users) TableName() string {
	return "users"
}

func TestSession_Execute(t *testing.T) {
	var sql = `CREATE TABLE IF NOT EXISTS "users" (
	 "uid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "name" TEXT NOT NULL,
	 "age" integer NOT NULL
)`
	var s = NewSession(initDB())
	var err error
	var aff int64

	aff, err = s.Execute(sql)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(aff)

	aff, err = s.Execute("insert into users(name,age) VALUES(?,?),(?,?)",
		"fizz", 18, "gorose", 19)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(aff)
}

func TestSession_Query_struct(t *testing.T) {
	var s = NewSession(initDB())
	var err error
	defer s.Close()

	var user []Users
	err = s.Table(&user).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("多条struct绑定:", user)

	var user2 Users
	err = s.Table(&user2).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("一条struct绑定:", user2)
}

//type UserMap map[string]interface{}

type aaa MapRow

func (u *aaa) TableName() string {
	return "users"
}

type bbb MapRows

func (u *bbb) TableName() string {
	return "users"
}

func TestSession_Query_map(t *testing.T) {
	var s = NewSession(initDB())
	var err error

	//var user = make([]map[string]interface{},0)
	var user = make(bbb, 0)
	err = s.Table(&user).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("多条map绑定:", user)

	var user2 = make(aaa)
	err = s.Table(&user2).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("一条map绑定:", user2)
}

func TestSession_Transaction(t *testing.T) {
	var s = NewSession(initDB())
	// 一键事务, 自动回滚和提交, 我们只需要关注业务即可
	err := s.Transaction(func1, func2)
	t.Log(err)
}

func func1(s ISession) error {
	var err error
	var aff int64
	aff, err = s.Execute("update users set name='?',age=? where uid=?",
		"gorose3", 21, 3)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("fail")
	}

	aff, err = s.Execute("update users set name='?',age=? where uid=?",
		"gorose2", 20, 2)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("fail")
	}

	return nil
}
func func2(s ISession) error {
	var err error
	var aff int64
	aff, err = s.Execute("update users set name='?',age=? where uid=?",
		"gorose3", 21, 3)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("fail")
	}

	aff, err = s.Execute("update users set name='?',age=? where uid=?",
		"gorose2", 20, 2)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("fail")
	}

	return nil
}
