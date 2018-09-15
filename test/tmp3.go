package main

import (
	"fmt"
	"reflect"
)

type StatusVal int
type Foo struct {
	Name string
	Age  int
}
type Bar struct {
	Status StatusVal
	FSlice []Foo
}

func ListFields(a interface{}) {
	v := reflect.ValueOf(a).Elem()
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		n := v.Type().Field(j).Name
		t := f.Type().String()
		fmt.Printf("Name: %s,  Kind: %s,  Type: %s\n", n, f.Kind(), t)
		fmt.Println(v.Type())
	}
}

func main() {
	var x Bar
	ListFields(&x)
}
