package gorose

import (
	"database/sql"
	"fmt"
	"github.com/gohouse/utils"
)


//var conn Connection
var DB *sql.DB
var Tx *sql.Tx
//var Config map[string]map[string]string

func Open(arg ...interface{}){
	var conn Connection
	if len(arg) == 1 {
		conn.Connect(arg[0])
	} else {
		conn.Config_ = arg[0].(map[string]map[string]string)
		conn.Connect(arg[1])
	}
}

type Connection struct {
	Config_ map[string]map[string]string
}

//func (this *Connection) Config(conf map[string]map[string]string) *Connection {
//	this.Config_ = conf
//
//	return this
//}

func (this *Connection) Connect(arg interface{}) *sql.DB {
	var err error
	var dbObj map[string]string

	if utils.GetType(arg) == "string" {
		dbObj = this.Config_[arg.(string)]
	} else {
		dbObj = arg.(map[string]string)
	}

	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s",
		dbObj["username"], dbObj["password"], dbObj["protocol"], dbObj["host"], dbObj["port"], dbObj["database"], dbObj["charset"])
	//DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	DB, err = sql.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}

	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		//log.Fatal(err.Error())
		panic(err.Error())
	}

	return DB
}

//func Config(conf map[string]map[string]string) *Connection {
//	var conn *Connection
//
//	res := conn.Config(conf)
//
//	fmt.Println(res)
//	return conn
//}
//func Connect(arg interface{}) *sql.DB {
//	var conn *Connection
//	conn.DB = conn.Connect(arg)
//
//	return conn.DB
//}

