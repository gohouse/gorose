package gorose

import (
	"testing"
)

func TestEngin(t *testing.T) {
	e := initDB()
	e.SetPrefix("pre_")

	t.Log(e.GetPrefix())

	db := e.GetQueryDB()

	err := db.Ping()

	if err!=nil {
		t.Error("gorose初始化失败")
	}
	t.Log("gorose初始化成功")
	t.Log(e.GetLogger())
}
