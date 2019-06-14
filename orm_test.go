package gorose

import (
	"fmt"
	"testing"
)

func initOrm() IOrm {
	return NewOrm(NewSession(initDB()), NewBinder())
}
func TestNewOrm(t *testing.T) {
	orm := initOrm()
	orm.Hello()
}

func TestOrm_Get(t *testing.T) {
	orm := initOrm()

	var u = aaa{}
	err := orm.Table(&u).Get()
	fmt.Println(u)
	fmt.Println(err)
}
