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
	ormObj := orm.Table(&u).Join("b", "a.id", "=", "b.id").
		Fields("uid,age").
		Order("uid desc").
		Where("a", 1).
		OrWhere(func() {
			orm.Where("c", 3).OrWhere(func() {
				orm.Where("d", ">", 4)
			})
		}).Where("e", 5).
		Limit(5).Offset(2)
	s, a, err := ormObj.BuildSql()

	fmt.Println(u)
	fmt.Println(err, s, a)
	fmt.Println(orm.LastSql())
}

func TestOrm_Pluck(t *testing.T) {
	orm := initOrm()

	var u = bbb{}
	ormObj := orm.Table(&u)

	//s, a, err := ormObj.BuildSql()
	//fmt.Println(u)
	//fmt.Println(err, s, a)
	//fmt.Println(orm.LastSql())

	//res,err := ormObj.Pluck("uid")
	err := ormObj.Get()
	fmt.Println(err)
	fmt.Println(u)
	//fmt.Println(res)
}

func TestOrm_Get2(t *testing.T) {
	orm := initOrm()

	var u = bbb{}
	ormObj := orm.Table(&u)

	res,err := ormObj.Pluck("name", "uid")
	//err := ormObj.Get()
	fmt.Println(err)
	fmt.Println(u)
	fmt.Println(res)
}
