package main

import "fmt"

func main() {
	var userMap map[int]string = map[int]string{1: "a", 2: "b", 3: "c"}
	fmt.Println(userMap)
	fmt.Println(userMap[1])

	delete(userMap, 2) //删除
	fmt.Println(userMap)

	//map 用之前要初始化
	//两种方法
	//var aMap = map[string]string{}
	//var aMap = make(map[string]string)
}
