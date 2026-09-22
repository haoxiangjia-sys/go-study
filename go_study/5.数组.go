package main

import "fmt"

func main() {
	var nameList [3]string = [3]string{"liu", "deng", "jia"}
	fmt.Println(nameList)
	fmt.Println(nameList[0])               //依靠索引 拿名字
	fmt.Println(len(nameList))             //长度
	fmt.Println(nameList[len(nameList)-1]) //最后一个

}
