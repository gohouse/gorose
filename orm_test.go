package gorose

import (
	"fmt"
	"testing"
)

func initOrm() IOrm {
	return NewOrm(NewSession(initDB()))
}
func TestNewOrm(t *testing.T) {
	orm := initOrm()
	orm.Hello()
}

func TestOrm_Get(t *testing.T) {
	orm := initOrm()

	var u = aaa{}
	err := orm.Table(&u).Join("b", "a.id", "=", "b.id").
		Fields("uid,age").
		Order("uid desc").
		Where("a", 1).
		OrWhere(func() {
			orm.Where("c", 3).OrWhere(func() {
				orm.Where("d", ">", 4)
			})
		}).Where("e", 5).
		Limit(5).Offset(2).Get()

	fmt.Println(u)
	fmt.Println(err)
	fmt.Println(orm.LastSql())
}
