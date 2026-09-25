package main

import "fmt"

func addOne(a int) int {
	return a + 1
}

func square(a int) int {
	return a * a
}

func double(slice []int) []int {
	slice = append(slice, slice...)
	//append(slice, elements...) 用于向切片末尾追加元素
}

func mapSlice(f func(a int) int, slice []int) {
	for index, value := range slice {
		mappedValue := f(value)
		slice[index] = mappedValue
	}
}

func mapArray(f func(a int) int, array [3]int) [3]int {
	var three [3]int
	for index, value := range array {
		three[index] = f(value)
	}
	return three
}

func main() {
	intsSlice := []int{1, 2, 3}
	mapSlice(addOne, intsSlice)
	fmt.Println(intsSlice)

	intsArray := [3]int{1, 2, 3}
	intsArray = mapArray(addOne, intsArray)
	fmt.Println(intsArray)

}
