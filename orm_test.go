package gorose

import (
	"fmt"
	"testing"
)

func TestNewOrm(t *testing.T) {
	orm := NewOrm(NewSession(initDB()), NewBinder())

	var u = MapRow{}
	err := orm.Table(&u).Get()
	fmt.Println(err)
}
