package gorose

import (
	"testing"
	"time"
)

func TestStructToMap(t *testing.T) {
	user := Users{Uid: 1, Name: "gorose"}
	data := StructToMap(user)
	t.Log(data)
}

func TestIf(t *testing.T) {
	closer := func() {
		time.Sleep(1 * time.Second)
	}
	withRunTimeContext(closer, func(td time.Duration) {
		t.Log("用时:", td, td.Seconds() > 1)
	})
}

//func TestStructToMap2(t *testing.T) {
//	var u Users
//	//res := structForScan(&u)
//	res := structForScan(reflect.ValueOf(&u).Interface())
//	for _, item := range res {
//		err := varBindValue.BindVal(item, 1234)
//		if err != nil {
//			t.Error(err.Error())
//		}
//	}
//	t.Log(res, u)
//}
func Test_getRandomInt(t *testing.T) {
	t.Log(getRandomInt(2))
	t.Log(getRandomInt(3))
	t.Log(getRandomInt(2))
	t.Log(getRandomInt(3))
}
