package main

import "github.com/gohouse/gorose"

type user3 struct {
	*gorose.Session
	Name string
}

func (u user3) TableName() string {
	return "users"
}

func main() {
	var user user3
	gorose.NewOrm().Table(&user).Select()
}
