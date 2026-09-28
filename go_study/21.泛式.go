package main

import "fmt"

// 定义一个类型参数 T，它可以是 int、float64 或 int32
// (a, b T) 是什么参数 a 的类型是 T  参数 b 的类型也是 T
// 表示函数返回值的类型也是 T
func add[T int | float64 | int32](a, b T) T {
	return a + b
}

func main() {
	fmt.Println(add(1, 2))               // T = int
	fmt.Println(add(1.5, 2.5))           // T = float64
	fmt.Println(add(int32(1), int32(2))) // T = int32
}

//1. 定义泛型结构体
//type Response[T any] struct {
//	Code int    `json:"code"`
//	Msg  string `json:"msg"`
//	Data T      `json:"data"`
//}

//泛型切片
//Go
//复制代码
//package main
//
//type MySlice[T any] []T
//
//func main() {
//  var mySlice MySlice[string]
//  mySlice = append(mySlice, "枫枫")
//  var intSlice MySlice[int]
//  intSlice = append(intSlice, 2)
//}
//泛型map
//Go
//复制代码
//package main
//
//import "fmt"
//
//type MyMap[K string | int, V any] map[K]V
//
//func main() {
//  var myMap = make(MyMap[string, string])
//  myMap["name"] = "枫枫"
//  fmt.Println(myMap)
//}
