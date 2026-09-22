package main

import "fmt"

func main() {
	var age int
	fmt.Printf("请输入你的名字： ")
	fmt.Scanln(&age)

	switch {
	case age < 18:

		fmt.Print("a")
	case age > 18:
		fmt.Print("b")
	default:
		fmt.Print("c")
	}
	//----------------------------------------------------
	var week int
	fmt.Printf("请输入星期： ")
	fmt.Scanln(&week)
	switch week {

	case 1:
		fmt.Print("a")
	case 2:
		fmt.Print("b")
	case 3:
		fmt.Print("c")
	default:
		fmt.Print("d")
	}
}
