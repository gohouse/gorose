package main

import (
	"fmt"
)

type Monkey struct {
	Name string
}


type BirdAble interface {
	Flying()
}

type Fishable interface {
	Swimming()
}

type LittleMonkey struct {
	Monkey //继承
}

func (this *Monkey) climbing() {
	fmt.Println(this.Name, "生来会爬树")
}

func (this *LittleMonkey) Flying() {
	fmt.Println(this.Name, "通过学习，我会飞翔了")
}

func (this *LittleMonkey) Swimming() {
	fmt.Println(this.Name, "通过学习，我会游泳了")
}

func main() {
	var _ BirdAble = &LittleMonkey{}
	var _ Fishable = &LittleMonkey{}
	monkey := LittleMonkey{
		Monkey{
			Name: "悟空",
		},
	}

	monkey.climbing()
	monkey.Flying()
	monkey.Swimming()
}

//当 A 结构体继承了 B 结构体，那么 A 结构就自动的继承了 B 结构体的字段和方法，并且可以直接使用
//当 A 结构体需要扩展功能，同时不希望去破坏继承关系，则可以去实现某个接口即可，因此我们可以认为:实现接口是对继承机制的补充.


//接口和继承解决的解决的问题不同
//继承的价值主要在于:解决代码的复用性和可维护性。
//接口的价值主要在于:设计，设计好各种规范(方法)，让其它自定义类型去实现这些方法。

//接口比继承更加灵活 Person Student BirdAble LittleMonkey
//接口比继承更加灵活，继承是满足 is - a 的关系，而接口只需满足 like - a 的关系。

//接口在一定程度上实现代码解耦