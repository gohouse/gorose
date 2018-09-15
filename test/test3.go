package main

import "fmt"

type a struct {

}
func main() {

	res := a{}
	res2 := new(a)

	fmt.Println(res, res2)
}
