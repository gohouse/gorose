package main

import (
	"github.com/gohouse/gorose/examples/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	"errors"
	"time"
)
var connection *gorose.Connection
var err error

func main() {
	connection = config.GetConnection()
	defer connection.Close()

	for k := 0; k < 5; k++ {
		go transTest(k)
	}
	time.Sleep(time.Second*1)
}

func transTest(k int)  {
	data2 := map[string]interface{}{
		"age":  17,
		"job":  "it31"+string(k),
		"name": "fizz4",
	}
	db := connection.NewSession()

	trans,err := db.Transaction(func() error {
		res2, err2 := db.Table("users").Data(data2).Insert()
		if err2 != nil {
			return err2
		}
		if res2 == 0 {
			return errors.New("Insert failed")
		}
		fmt.Println(res2)
		return nil
	})

	fmt.Println(trans, err)
}
