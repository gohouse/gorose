package gorose

import (
	"fmt"
	"testing"
)

func initOrm() IOrm {
	return initDB().NewOrm()
}
func DB() IOrm {
	return initDB().NewOrm()
}
func TestNewOrm(t *testing.T) {
	orm := initOrm()
	orm.Hello()
}
func TestOrm_AddFields(t *testing.T) {
	orm := initOrm()
	var u = User{}
	var fieldStmt = orm.Table(&u).Fields("a").Where("m",55)
	a,b,c := fieldStmt.AddFields("b").Where("d",1).BuildSql()
	fmt.Println(a,b,c)

	d,e,f := fieldStmt.AddFields("c").Where("d",2).BuildSql()
	fmt.Println(d,e,f)
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

	//var u = aaa{}
	var u = bbb{}
	//var u Users
	//var u []Users
	ormObj := orm.Table(&u)
	//res,err := ormObj.Pluck("name", "uid")
	res, err := ormObj.Pluck("name")
	fmt.Println(err)
	fmt.Println(u)
	fmt.Println(res)
}

func TestOrm_Value(t *testing.T) {
	orm := initOrm()

	//var u = aaa{}
	var u = bbb{}
	//var u Users
	//var u []Users
	ormObj := orm.Table(&u)
	res, err := ormObj.Value("name")
	fmt.Println(err)
	fmt.Println(u)
	fmt.Println(res)
}

func TestOrm_Count(t *testing.T) {
	orm := initOrm()

	var u = aaa{}
	ormObj := orm.Table(&u)
	res, err := ormObj.Count()
	fmt.Println(err)
	fmt.Println(orm.LastSql())
	fmt.Println(res)
}

func TestOrm_Chunk(t *testing.T) {
	orm := initOrm()

	var u = bbb{}
	err := orm.Table(&u).Chunk(1, func(data []Map) error {
		for _,item := range data{
			fmt.Println(item["name"].String())
		}
		return nil
	})
	fmt.Println(err)
}

func TestOrm_Loop(t *testing.T) {
	db := DB()

	var u = bbb{}
	//aff,err := db.Table(&u).Force().Data(Data{"age": 18}).Update()
	//fmt.Println(aff,err)
	err := db.Table(&u).Where("age", 18).Loop(2, func(data []Map) error {
		for _,item := range data{
			DB().Data(Data{"age": 19}).Where("uid", item["uid"].Int64()).Update()
		}
		return nil
	})
	fmt.Println(err)
}

func TestOrm_Get2(t *testing.T) {
	db := initOrm()
	var u = bbb{}
	var err error

	err = db.Table(&u).Limit(1).Get()
	fmt.Println(err,u,db.LastSql())

	res,err := db.Table("users").Limit(2).GetAll()
	fmt.Println(err,res,db.LastSql())

	fmt.Println(u)
}
