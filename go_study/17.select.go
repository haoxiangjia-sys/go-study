//
//
//select {
//case v := <-ch1:
//fmt.Println("来自 ch1:", v)
//case v := <-ch2:
//fmt.Println("来自 ch2:", v)
//case ch3 <- 100:
//fmt.Println("往 ch3 发送成功")
//default:
//fmt.Println("都没有准备好")
//}//它让一个 goroutine 同时等待多个 channel 操作，谁先准备好就先处理谁。

package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	done := make(chan struct{})

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "来自 ch1"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "来自 ch2"
	}()
	go func() {
		time.Sleep(3 * time.Second)
		close(done)
	}()

	for {
		select {
		case v := <-ch1:
			fmt.Println(v)
		case v := <-ch2:
			fmt.Println(v)
		case <-done:
			fmt.Println("全部结束")
			return
		case <-time.After(5 * time.Second):
			fmt.Println("总超时")
			return
		}
	}
}
