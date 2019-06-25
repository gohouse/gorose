package gorose

import (
	"fmt"
	"testing"
)

func TestOrm_First(t *testing.T) {
	db := initOrm()
	var u = Users{}
	res, err := db.Table(&u).Get()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(res)
	db.Table(&u).Limit(2).Select()
	t.Log(u)
}

func TestOrm_Select(t *testing.T) {
	db := initOrm()
	var err error

	var u = []Users{}
	err = db.Table(&u).Limit(2).Select()
	fmt.Println(err, u, db.LastSql())

	var u2 = Users{}
	err = db.Table(&u2).Limit(1).Select()
	fmt.Println(err, u2, db.LastSql())

	var u3 Users
	err = db.Table(&u3).Limit(1).Select()
	fmt.Println(err, u3, db.LastSql())

	var u4 []Users
	err = db.Table(&u4).Limit(2).Select()
	fmt.Println(err, u4, db.LastSql())
}

func TestOrm_Select2(t *testing.T) {
	db := initOrm()
	var err error

	var u = UsersMap{}
	err = db.Table(&u).Limit(2).Select()
	fmt.Println(err, u, db.LastSql())

	var u3 = UsersMapSlice{}
	err = db.Table(&u3).Limit(1).Select()
	fmt.Println(err, u3, db.LastSql())
}

func TestOrm_Get2(t *testing.T) {
	db := initOrm()
	var err error
	var u = []Users{}

	//res, err := db.Table("users").Where("uid", ">", 2).Limit(2).Get()
	res, err := db.Table(&u).Where("uid", ">", 2).Limit(2).Get()
	fmt.Println(err, res, db.LastSql())

	fmt.Println(u)
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
		for _, item := range data {
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
		fmt.Println(db.LastSql())
		for _, item := range data {
			DB().Table(&u).Data(Data{"age": 19}).Where("uid", item["uid"].Int64()).Update()
		}
		return nil
	})
	fmt.Println(err)
}
