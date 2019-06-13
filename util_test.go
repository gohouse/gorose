package gorose

import (
	"testing"
	"time"
)

type User struct {
	Id        int64
	Username  string
	Password  string
	LoginTime time.Time
}

func TestStructToMap(t *testing.T) {
	user := User{5, "Wall", "pwd", time.Now()}
	data := StructToMap(user)
	t.Log(data)
}
