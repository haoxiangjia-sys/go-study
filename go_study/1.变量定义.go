package main

import "fmt"

const version = "2.0.2" //常量 定义后无法修改

func main() {
	fmt.Printf("Version: %s\n", version)
	//声明
	var name string
	//赋值
	name = "xiang"
	fmt.Println(name)

	//声明 + 赋值
	var name1 string = "xiang"
	fmt.Println(name1)
	//省略
	var name2 = "xiang"
	fmt.Println(name2)
	//短声明
	name3 := "xiang"
	fmt.Println(name3)

	fmt.Println(version)
}

/*run
1,终端 go run main.go
2, 点空白区域  run

全局变量
局部变量

命名规范 ： 首字母大写的变量 函数 方法 属性 可以在包外访问
*/
