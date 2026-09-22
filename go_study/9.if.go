package main

import "fmt"

func main() {
	var age int
	fmt.Printf("请输入你的名字： ")
	fmt.Scanln(&age)

	/*
		//中断式
		if age <= 0 {
			fmt.Println("未出生")
			return
		}
		if age <= 18 {
			fmt.Println("未成年")
			return
		}
		fmt.Println("中年")
	*/

	/*
		//嵌入式
		if age > 18 {
			if age <= 0 {
				fmt.Println("未出生")
			} else {
				fmt.Println("未成年")
			}
		} else {
			fmt.Println("中年")
		}
	*/

	/*
		if age <= 0 {
			fmt.Println("未出生")
		}
		if age <=18 && age > 0{
			fmt.Print()
		}

	*/

	// && and  ; \\ or
	//&& 第一个条件是false 停
	//|| 第一个     true  停
}
