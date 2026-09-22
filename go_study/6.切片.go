package main

import (
	"fmt"
	"sort"
)

func main() {
	var nameList []string //slice 比 数组更加灵活 //数组 长度改不了
	nameList = append(nameList, "liu")
	nameList = append(nameList, "jia")
	nameList = append(nameList, "chang")
	fmt.Println(nameList)

	//make 函数
	//make([]type, length, capacity)
	ageList := make([]int, 3)
	fmt.Println(ageList)

	//切
	array := [3]int{1, 2, 3}
	fmt.Println(array[0:2]) //从二开始切两位
	fmt.Println(array[1:2]) //从一开始切一位

	var ints = []int{4, 6, 3, 5}
	sort.Ints(ints) //升序
	fmt.Println(ints)
	sort.Sort(sort.Reverse(sort.IntSlice(ints)))
	fmt.Println(ints)

}
