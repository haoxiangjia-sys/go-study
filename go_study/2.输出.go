package main

import "fmt"

func main() {
	fmt.Println("hello", "nihao")
	fmt.Println("nihao", "xiangxiang")

	fmt.Printf("%s hello \n", "浩翔") //\n 是换行符（newline）—— 让输出从这里开始换到下一行。
	fmt.Printf("%.2f\n", 3.141592)
	fmt.Printf("%T %T\n", "你好", 2.5) //%T 查找类型
	fmt.Printf("%#v\n", "")          //用go 语言格式输出，打印空字符串
}
