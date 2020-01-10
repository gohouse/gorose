package gorose

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestOrm_BuildSql2(t *testing.T) {
	db := DB()
	var u = "age=age+1,num=num+1"
	var wheres interface{}
	wheres = [][]interface{}{{"a", ">", "b"}, {"a", "b"},{"a is null"}}
	sqlstr, a, b := db.Force().Table("users").Data(u).Where(wheres).BuildSql("update")

	t.Log(sqlstr, a, b)
}

func TestOrm_BuildSql3(t *testing.T) {
	db := DB()
	var u = "age=age+1,num=num+1"
	var wheres interface{}
	wheres = [][]interface{}{{"a", ">", "b"}, {"a", "b"}}
	sqlstr, a, b := db.Force().Table(Users{}).Data(u).Where(wheres).BuildSql("update")

	t.Log(sqlstr, a, b)
}

func TestOrm_BuildSql4(t *testing.T) {
	//sqlstr, a, b := db.Table("users3").Limit(2).Offset(2).BuildSql()
	//t.Log(sqlstr, a, b)
	//
	//sqlstr, a, b = db.Table("users2").Limit(2).Offset(2).BuildSql()
	//t.Log(sqlstr, a, b)

	//var u = Users{
	//	Uid:  1,
	//	Name: "2",
	//	Age:  3,
	//}
	//res,err := db.Table("xxx").Insert(&u)
	//t.Log(db.LastSql(), res,err)
}

func TestOrm_BuildSql5(t *testing.T) {
	//ticker := time.NewTicker(100*time.Millisecond)
	go func() {
		for {
			//<-ticker.C
			db := DB()
			sqlstr, a, b := db.Table("users").Where("uid", ">", 1).BuildSql()
			//c,d := db.Table("users").Get()
			//t.Log(db.LastSql())
			count, d := db.First()

			t.Log(sqlstr, a, b)
			t.Log(count, d)
			t.Log(db.LastSql())
		}
	}()
	time.Sleep(500 * time.Millisecond)
}

func TestOrm_BuildSql6(t *testing.T) {
	var db = DB()
	sqlstr, a, b := db.Table("users3").Limit(2).Offset(2).BuildSql()
	t.Log(sqlstr, a, b)

	sqlstr, a, b = db.Table("users2").Limit(2).Offset(2).BuildSql()
	t.Log(sqlstr, a, b)

	var u = Users{
		Uid:  1111,
		Name: "2",
		Age:  3,
	}
	res,err := db.Table("xxx").Where("xx","xx").Update(&u)
	t.Log(db.LastSql(), res,err)
}

func TestOrm_First(t *testing.T) {
	res, err := DB().Table(Users{}).Where("uid",1).First()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res)
}

func TestOrm_Select(t *testing.T) {
	db := DB()
	var err error

	var u = []Users{}
	err = db.Table(&u).Select()
	t.Log(err, u, db.LastSql())

	var u2 = Users{}
	err = db.Table(&u2).Select()
	t.Log(err, u2, db.LastSql())

	var u3 Users
	err = db.Table(&u3).Select()
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

	var u = []UsersMap{}
	err = db.Table(&u).Limit(2).Select()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(u)

	var u3 = UsersMap{}
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

	res, err := db.Table("users").Where("uid", ">", 2).
		//Where("1","=","1").
		Where("1 = 1").
		Limit(2).Get()
	//res, err := db.Table(&u).Where("uid", ">", 0).Limit(2).Get()
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
		Where("a", "like", "%3%").
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
	//var u []Users
	ormObj := orm.Table("users")
	//res,err := ormObj.Pluck("name", "uid")
	res, err := ormObj.Limit(5).Pluck("name","uid")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, orm.LastSql())
}

func TestOrm_Value(t *testing.T) {
	db := DB()

	//var u = UsersMap{}
	//var u = UsersMapSlice{}
	//var u Users
	//var u []Users
	//ormObj := db.Table(&u)
	//ormObj := db.Table("users")
	ormObj := db.Table(Users{})
	res, err := ormObj.Limit(5).Value("uid")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, db.LastSql())
}

func TestOrm_Count(t *testing.T) {
	db := DB()

	//var u = UsersMap{}
	//ormObj := db.Table(&u)
	ormObj := db.Table("users")

	res, err := ormObj.Count()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, db.LastSql())
}

func TestOrm_Count2(t *testing.T) {
	var u Users
	var count int64
	count, err := DB().Table(&u).Count()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(count)
}

func TestOrm_Chunk(t *testing.T) {
	orm := DB()

	var u = []UsersMap{}
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

func TestOrm_Chunk2(t *testing.T) {
	orm := DB()

	var u []Users
	var i int
	err := orm.Table(&u).ChunkStruct(2, func() error {
		//for _, item := range u {
		t.Log(u)
		//}
		if i == 2 {
			return errors.New("故意停止,防止数据过多,浪费时间")
		}
		i++
		return nil
	})
	if err != nil && err.Error() != "故意停止,防止数据过多,浪费时间" {
		t.Error(err.Error())
	}
	t.Log("ChunkStruct() success")
}

func TestOrm_Loop(t *testing.T) {
	db := DB()

	var u = []UsersMap{}
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
	res, err := db.Table(&u).Limit(2).Paginate()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, u)
	t.Log(db.LastSql())
}

func TestOrm_Paginate2(t *testing.T) {
	db := DB()

	var u []Users
	res, err := db.Table(&u).Where("uid",">",1).Limit(2).Paginate(3)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res, u)
	t.Log(db.LastSql())
}

func TestOrm_Sum(t *testing.T) {
	db := DB()

	var u Users
	//res, err := db.Table(Users{}).First()
	res, err := db.Table(&u).Where(Data{"uid":1}).Sum("age")
	if err != nil {
		t.Error(err.Error())
	}
	//fmt.Printf("%#v\n",res)
	t.Log(res, u)
	t.Log(db.LastSql())
}

func BenchmarkNewOrm(b *testing.B) {
	engin := initDB()
	for i := 0; i < b.N; i++ {
		engin.NewOrm().Table("users").First()
	}
}

func BenchmarkNewOrm2(b *testing.B) {
	engin := initDB()
	for i := 0; i < b.N; i++ {
		engin.NewOrm().Table("users").First()
	}
}
