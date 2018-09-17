package main

import (
	"fmt"
	"github.com/gohouse/gorose"
)

func main() {
	var a gorose.Config
	fmt.Println(gorose.Config{}.DbConfig)
}
