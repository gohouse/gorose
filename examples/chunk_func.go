package main

import (
	"fmt"
)

//var data = []int {1,2,3,4,5,6,7,8,9,0}
var Data []int

func init()  {
	Data = append(Data, 1,2,3,4,5,6,7,8,9,0)
}

func main() {
	chunk(2, func(result []int) {
		fmt.Println(result)
	})
}


func chunk(limit int, callback func(arg []int))  {
	var step = 0
	for {
		start := step*limit
		if len(Data)<=start {
			break
		}
		current_data := Data[start:limit+start]

		var res []int
		for _, item := range current_data {
			res = append(res, item)
		}

		callback(res)
//fmt.Println(res)
		if len(current_data)<limit {
			break
		}
		step++
	}
}
