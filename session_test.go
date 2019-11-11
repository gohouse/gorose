package gorose

import (
	"errors"
	"testing"
)

func initSession() ISession {
	return initDB().NewSession()
}

func TestSession_Query(t *testing.T) {
	var s = initSession()
	//var user []Users
	res, err := s.Query("select * from users where name=?", "gorose2")
	if err != nil {
		t.Error(err.Error())
	}
	//t.Log(res)
	t.Log(res, s.LastSql())
}

func TestSession_Query3(t *testing.T) {
	var s = initSession()
	var o []Users
	//var o []map[string]interface{}
	//var o []gorose.Data
	res, err := s.Bind(&o).Query("select * from users limit 2")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res)
	t.Log(o, s.LastSql())
}

func TestSession_Execute(t *testing.T) {
	var sql = `CREATE TABLE IF NOT EXISTS "orders" (
	 "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "goodsname" TEXT NOT NULL default "",
	 "price" decimal default "0.00"
	)`
	var s = initSession()
	var err error
	var aff int64

	aff, err = s.Execute(sql)
	if err != nil {
		t.Error(err.Error())
	}
	if aff == 0 {
		return
	}

	aff, err = s.Execute("insert into orders(goodsname,price) VALUES(?,?),(?,?)",
		"goods1", 1.23, "goods2", 3.23)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(aff)
}

func TestSession_Query_struct(t *testing.T) {
	var s = initSession()
	var err error
	// defer s.Close()

	var user []Users
	_, err = s.Bind(&user).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("多条struct绑定:", user)

	var user2 Users
	_, err = s.Bind(&user2).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("一条struct绑定:", user2)
}

//type UserMap map[string]interface{}

func TestSession_Query_map(t *testing.T) {
	var s = initSession()
	var err error

	var user2 = aaa{}
	_, err = s.Bind(&user2).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("一条map绑定:", user2)
	t.Log("一条map绑定的uid为:", user2["uid"])
	t.Log(s.LastSql())

	var user = bbb{}
	_, err = s.Bind(&user).Query("select * from users limit ?", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("多条map绑定:", user)
	t.Log("多条map绑定:", user[0]["age"].Int())
	t.Log(s.LastSql())
}

func TestSession_Bind(t *testing.T) {
	var s = initSession()
	var err error

	var user2 = aaa{}
	_, err = s.Bind(&user2).Query("select * from users limit ?", 2)

	if err != nil {
		t.Error(err.Error())
	}
	t.Log("session bind success")
}

func TestSession_Transaction(t *testing.T) {
	var s = initSession()
	// 一键事务, 自动回滚和提交, 我们只需要关注业务即可
	err := s.Transaction(trans1, trans2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("session transaction success")
}

func trans1(s ISession) error {
	var err error
	var aff int64
	aff, err = s.Execute("update users set name=?,age=? where uid=?",
		"gorose3", 21, 3)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("fail")
	}

	aff, err = s.Execute("update users set name=?,age=? where uid=?",
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
	aff, err = s.Execute("update users set name=?,age=? where uid=?",
		"gorose3", 21, 3)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("fail")
	}

	aff, err = s.Execute("update users set name=?,age=? where uid=?",
		"gorose2", 20, 2)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("fail")
	}

	return nil
}
