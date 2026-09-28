package main

import (
	"fmt"
	"runtime/debug"
)

////向上抛：把错误交给上一级
//package main
//
//import (
//"errors"
//"fmt"
//)
//
//func Parent() error {
//	err := method()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func method() error {
//	return errors.New("出错了")
//}
//
//func main() {
//	fmt.Println(Parent())
//}

//中断程序：遇到错误直接停止
//package main
//
//import (
//"fmt"
//"os"
//)
//
//func init() {
//	_, err := os.ReadFile("xxx")
//	if err != nil {
//		panic(err)
//	}
//}
//
//func main() {
//	fmt.Println("啦啦啦")
//}


恢复程序：捕获 panic，防止崩溃
含义
用 defer + recover 捕获 panic，让程序不至于直接崩溃。

示例
go
package main

import (
"fmt"
"runtime/debug"
)

func read() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			s := string(debug.Stack())
			fmt.Println(s)
		}
	}()

	var list = []int{2, 3}
	fmt.Println(list[2]) // 越界，触发 panic
}

func main() {
	read()
	fmt.Println("程序继续执行")
}