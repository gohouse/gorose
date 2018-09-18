package main

import (
	"fmt"
	"github.com/gohouse/gorose"
	"reflect"
)


func main() {
	var conf gorose.Configure
	//var conf2 across.DbConfigSingle

	conf.Prefix = "pre_"

	fmt.Println(conf)
	fmt.Println(reflect.TypeOf(conf.DbConfigSingle))
}
