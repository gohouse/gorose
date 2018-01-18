package main

import (
	"fmt"
)

func main() {
	a := []interface{}{"a","b"}
fmt.Println(a...)
	tt(a...)
}

func tt(args ...interface{}) {
	fmt.Println(args)
}