package main

import (
	"fmt"
)

// Data tmp
var Data []int

func init() {
	Data = append(Data, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0)
}

func main() {
	chunk(2, func(result []int) {
		fmt.Println(result)
	})
}

func chunk(limit int, callback func(arg []int)) {
	var step = 0
	for {
		start := step * limit
		if len(Data) <= start {
			break
		}
		currentData := Data[start : limit+start]

		var res []int
		for _, item := range currentData {
			res = append(res, item)
		}

		callback(res)
		//fmt.Println(res)
		if len(currentData) < limit {
			break
		}
		step++
	}
}
