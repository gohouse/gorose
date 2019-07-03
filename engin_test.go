package gorose

import (
	"github.com/gohouse/t"
	"testing"
)

type aaa t.MapString

func (u *aaa) TableName() string {
	return "users"
}

//type bbb MapRows
type bbb []t.MapString

func (u *bbb) TableName() string {
	return "users"
}

type UsersMap Data

func (u *UsersMap) TableName() string {
	return "users"
}

// 定义map多返回绑定表名,一定要像下边这样,单独定义,否则无法获取对应的 TableName()
type UsersMapSlice []Data

func (u *UsersMapSlice) TableName() string {
	return "users"
}

type Users struct {
	Uid  int    `orm:"uid"`
	Name string `orm:"name"`
	Age  int    `orm:"age"`
	Fi   string `orm:"ignore"`
}

func (u *Users) TableName() string {
	return "users"
}

func TestEngin(t *testing.T) {
	e := initDB()
	e.SetPrefix("pre_")

	t.Log(e.GetPrefix())

	db := e.GetQueryDB()

	err := db.Ping()

	if err != nil {
		t.Error("gorose初始化失败")
	}
	t.Log("gorose初始化成功")
	t.Log(e.GetLogger())
}
