package gorose

import (
	"errors"
	"fmt"
	"github.com/gohouse/t"
	"testing"
)

type Users struct {
	Uid  int    `gorose:"uid"`
	Name string `gorose:"name"`
	Age  int    `gorose:"age"`
	Fi string `gorose:"fi"`
}

func (u *Users) TableName() string {
	return "users"
}

func initSession() ISession {
	return initDB().NewSession()
}

func TestSession_Query(t *testing.T) {
	var s = initSession()
	var user []Users
	err := s.Bind(&user).Query("select * from users where name=?", "gorose")
	fmt.Println(user, err)
	fmt.Println(s.LastSql())
}

func TestSession_Execute(t *testing.T) {
	var sql = `CREATE TABLE IF NOT EXISTS "users" (
	 "uid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "name" TEXT NOT NULL,
	 "age" integer NOT NULL
)`
	var s = initSession()
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
	var s = initSession()
	var err error
	defer s.Close()

	var user []Users
	err = s.Bind(&user).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("多条struct绑定:", user)

	var user2 Users
	err = s.Bind(&user2).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("一条struct绑定:", user2)
}

//type UserMap map[string]interface{}

type aaa t.MapString

func (u *aaa) TableName() string {
	return "users"
}

//type bbb MapRows
type bbb []t.MapString

func (u *bbb) TableName() string {
	return "users"
}

func TestSession_Query_map(t *testing.T) {
	var s = initSession()
	var err error

	var user2 = aaa{}
	err = s.Bind(&user2).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("一条map绑定:", user2)
	t.Log("一条map绑定的uid为:", user2["uid"])
	t.Log(s.LastSql())

	var user = bbb{}
	err = s.Bind(&user).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("多条map绑定:", user)
	t.Log("多条map绑定:", user[0]["age"].Int())
	t.Log(s.LastSql())
}

func TestSession_Orm(t *testing.T) {
	var s = initSession()
	var err error

	var user2 = aaa{}
	err = s.Bind(&user2).Query("select * from users limit ?", 2)

	fmt.Println(err)
}

func TestSession_Transaction(t *testing.T) {
	var s = initSession()
	// 一键事务, 自动回滚和提交, 我们只需要关注业务即可
	err := s.Transaction(trans1, trans2)
	t.Log(err)
}

func trans1(s ISession) error {
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
func trans2(s ISession) error {
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
