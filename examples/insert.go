package main

import (
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
	"time"
)

var engin *gorose.Engin
var err error
func init() {
	engin, err = gorose.Open(&gorose.Config{Prefix: "nv_", Driver: "mysql", Dsn: "root:123456@tcp(localhost:3306)/novel?charset=utf8&parseTime=true"})
	if err!=nil {
		panic(err.Error())
	}
}
func db() gorose.IOrm {
	return engin.NewOrm()
}
type Tag struct {
	Id        int64     `gorose:"id" json:"id"`
	TagTitle  string    `gorose:"tag_title" json:"tag_title"`
	TagStatus int64     `gorose:"tag_status" json:"tag_status"` // 标签状态：默认0未启用，1启用
	CreatedAt time.Time `gorose:"created_at" json:"created_at"`
}

func (Tag) TableName() string {
	return "tag"
}
func main() {
	//fmt.Println(strings.TrimPrefix("nv_novel_chapter", "nv_"))
	//return
	var t = Tag{
		TagTitle:  "0",
		TagStatus: 1,
	}
	res, err := db().Insert(&t)
	fmt.Println(res, err)
}
