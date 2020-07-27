package gorose

import (
	"github.com/gohouse/t"
	"testing"
)

type aaa t.MapStringT

func (u *aaa) TableName() string {
	return "users"
}

//type bbb MapRows
type bbb []t.MapStringT

func (u *bbb) TableName() string {
	return "users"
}

type UsersMap Data

func (*UsersMap) TableName() string {
	return "users"
}

// 定义map多返回绑定表名,一定要像下边这样,单独定义,否则无法获取对应的 TableName()
type UsersMapSlice []Data

func (u *UsersMapSlice) TableName() string {
	return "users"
}

type Users struct {
	Uid  int64  `orm:"uid"`
	Name string `orm:"name"`
	Age  int64  `orm:"age"`
	Fi   string `orm:"ignore"`
}

func (Users) TableName() string {
	return "users"
}

type Orders struct {
	Id        int     `orm:"id"`
	GoodsName string  `orm:"goodsname"`
	Price     float64 `orm:"price"`
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
