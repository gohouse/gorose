package gorose

import (
	"errors"
	"fmt"
	"testing"
)

func TestOrm_First(t *testing.T) {
	db := DB()
	var u = Users{}
	var err error
	//res, err := db.Table(&u).Get()
	//if err != nil {
	//	t.Error(err.Error())
	//}
	//t.Log(res)
	err = db.Table(&u).Select()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(u)
}

func TestOrm_Select(t *testing.T) {
	db := DB()
	var err error

	var u = []Users{}
	err = db.Table(&u).Limit(2).Select()
	t.Log(err, u, db.LastSql())

	var u2 = Users{}
	err = db.Table(&u2).Limit(1).Select()
	t.Log(err, u2, db.LastSql())

	var u3 Users
	err = db.Table(&u3).Limit(1).Select()
	t.Log(err, u3, db.LastSql())

	var u4 []Users
	err = db.Table(&u4).Limit(2).Select()
	t.Log(err, u4, db.LastSql())
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(u, u2, u3, u4)
}

func TestOrm_Select2(t *testing.T) {
	db := DB()
	var err error

	var u = UsersMap{}
	err = db.Table(&u).Limit(2).Select()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(u)

	var u3 = UsersMapSlice{}
	err = db.Table(&u3).Limit(1).Select()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(u)
}


type Users2 struct {
	Name string `orm:"name"`
	Age  int    `orm:"age"`
	Uid  int    `orm:"uid"`
	Fi   string `orm:"ignore"`
}

func (u *Users2) TableName() string {
	return "users"
}
func TestOrm_Get2(t *testing.T) {
	db := DB()
	var err error
	var u []Users2

	//res, err := db.Table("users").Where("uid", ">", 2).Limit(2).Get()
	res, err := db.Table(&u).Where("uid", ">", 0).Limit(2).Get()
	fmt.Println(db.LastSql())
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, u)
}

func TestOrm_Get(t *testing.T) {
	orm := DB()

	var u = UsersMap{}
	ormObj := orm.Table(&u).Join("b", "a.id", "=", "b.id").
		RightJoin("userinfo d on a.id=d.id").
		Fields("a.uid,a.age").
		Order("uid desc").
		Where("a", 1).
		WhereNull("bb").
		WhereNotNull("cc").
		WhereIn("dd", []interface{}{1, 2}).
		OrWhereNotIn("ee", []interface{}{1, 2}).
		WhereBetween("ff", []interface{}{11, 21}).
		WhereNotBetween("ff", []interface{}{1, 2}).
		OrWhere(func() {
			orm.Where("c", 3).OrWhere(func() {
				orm.Where("d", ">", 4)
			})
		}).Where("e", 5).Limit(5).Offset(2)
	s, a, err := ormObj.BuildSql()

	if err != nil {
		t.Error(err.Error())
	}
	t.Log(s, a, u)
}

func TestOrm_Pluck(t *testing.T) {
	orm := DB()

	//var u = UsersMapSlice{}
	var u []Users
	ormObj := orm.Table(&u)
	res,err := ormObj.Pluck("name", "uid")
	//res, err := ormObj.Limit(5).Pluck("name")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, orm.LastSql())
}

func TestOrm_Value(t *testing.T) {
	db := DB()

	//var u = UsersMap{}
	var u = UsersMapSlice{}
	//var u Users
	//var u []Users
	ormObj := db.Table(&u)
	res, err := ormObj.Limit(5).Value("name")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, db.LastSql())
}

func TestOrm_Count(t *testing.T) {
	db := DB()

	var u = UsersMap{}
	ormObj := db.Table(&u)
	res, err := ormObj.Count()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, db.LastSql())
}

func TestOrm_Chunk(t *testing.T) {
	orm := DB()

	var u = UsersMapSlice{}
	err := orm.Table(&u).Chunk(1, func(data []Data) error {
		for _, item := range data {
			t.Log(item["name"])
		}
		return errors.New("故意停止,防止数据过多,浪费时间")
		//return nil
	})
	if err != nil && err.Error() != "故意停止,防止数据过多,浪费时间" {
		t.Error(err.Error())
	}
	t.Log("Chunk() success")
}

func TestOrm_Loop(t *testing.T) {
	db := DB()

	var u = UsersMapSlice{}
	//aff,err := db.Table(&u).Force().Data(Data{"age": 18}).Update()
	//fmt.Println(aff,err)
	err := db.Table(&u).Where("age", 18).Loop(2, func(data []Data) error {
		for _, item := range data {
			_, err := DB().Table(&u).Data(Data{"age": 19}).Where("uid", item["uid"]).Update()
			if err != nil {
				t.Error(err.Error())
			}
		}
		return errors.New("故意停止,防止数据过多,浪费时间")
		//return nil
	})
	if err != nil && err.Error() != "故意停止,防止数据过多,浪费时间" {
		t.Error(err.Error())
	}
	t.Log("Loop() success")
}


func TestOrm_Paginate(t *testing.T) {
	db := DB()

	var u []Users
	res,err := db.Table(&u).Limit(2).Paginate()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, u)
	t.Log(db.LastSql())
}

func BenchmarkNewOrm(b *testing.B) {
	engin:=initDB()
	for i:=0;i<b.N;i++{
		engin.NewOrm().Table("users").First()
	}
}

func BenchmarkNewOrm2(b *testing.B) {
	engin:=initDB()
	for i:=0;i<b.N;i++{
		engin.NewOrm().Table("users").First()
	}
}
