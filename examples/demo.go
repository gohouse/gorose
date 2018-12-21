package main

import (
	"fmt"
	"github.com/gohouse/gorose"
	_ "github.com/gohouse/gorose/driver/mysql"
)

func main() {
	fmt.Println(gorose.VERSION)
}
