package main

import (
	"fmt"
	"sync"
)

//func main() {
//	fmt.Println("hello world")
//
//	//question1
//	for i := 0; i < 20; i++ {
//		fmt.Println("hello world")
//	}
//}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1) // 每启动一个 goroutine，计数器加 1
		go func(i int) {
			defer wg.Done() // goroutine 结束时计数器减 1
			fmt.Println("Hello from goroutine", i)
		}(i) // 把 i 作为参数传入，避免闭包捕获问题
	}

	wg.Wait() // 阻塞 main，直到所有 goroutine 完成
}
