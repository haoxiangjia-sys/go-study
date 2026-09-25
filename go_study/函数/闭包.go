package main

import (
	"fmt"
	"time"
)

// 设计一个函数， 先传一个参数表示延时，后面再传参数，将参数和
func awaitAdd(t int) func(...int) int {
	time.Sleep(time.Duration(t) * time.Second) ///****
	return func(numList ...int) int {
		var sum int
		for _, i2 := range numList {
			sum += i2
		}
		return sum
	}
}

func main() {
	fmt.Println(awaitAdd(10)(1, 2, 3, 4, 5))
}
