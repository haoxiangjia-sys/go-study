package main

import "fmt"

func sayHello() {
	fmt.Println("hello")
}

func param1(id string) {
	fmt.Println(id)
}

func param2(id, userName string) {
	fmt.Println(id, userName)
}

func add(numberList ...int) {
	var sum int
	for _, number := range numberList {
		sum += number
	}
	fmt.Println(sum)
}

func main() {
	sayHello()
	param1("123")
	param2("123", "JIA")
	add(1, 2, 3)
	add(1, 2, 3, 4, 5)

	var getName = func() string {
		return "xiangxiang"
	}
	var setName = func(name string) {
		fmt.Println(name)
		return
	}
	fmt.Println(getName())
	setName("kang")
}
