package main

import (
	"fmt"
	"reflect"
)

type YourT1 struct {
}

func (y *YourT1) MethodBar() {
	fmt.Println("MethodBar called")
}

type YourT2 struct {
}

func (y *YourT2) MethodFoo(i int, oo string) {
	fmt.Println("MethodFoo called", i, oo)
}

func InvokeObjectMethod(object interface{}, methodName string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(object).MethodByName(methodName).Call(inputs)
}

func main() {
	InvokeObjectMethod(new(YourT2), "MethodFoo", 10, "abc")
	InvokeObjectMethod(new(YourT1), "MethodBar")
}