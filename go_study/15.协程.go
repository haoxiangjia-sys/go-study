package main

import (
	"fmt"
	"sync"
	"time"
)

// // 声明一个全局的 WaitGroup 变量 wait
// // WaitGroup 用于等待一组 goroutine 执行完毕
var (
	wait = sync.WaitGroup{}
)

// // sing 模拟一个耗时 1 秒的任务
func sing() {
	fmt.Println("唱歌")
	time.Sleep(1 * time.Second)
	fmt.Println("唱歌结束")

	// 关键：通知 WaitGroup 当前 goroutine 已完成
	// 每调用一次 Done()，WaitGroup 内部的计数器就减 1
	wait.Done()
}

func main() {
	// 设置 WaitGroup 的计数器为 4
	// 表示接下来需要等待 4 个任务完成
	// 这个 Add 必须在启动 goroutine 之前调用
	wait.Add(4) // 计数器 = 4

	// 启动 4 个 goroutine 并发执行 sing()
	// 每个 goroutine 都是独立运行的，执行顺序不确定
	go sing()
	go sing()
	go sing()
	go sing()

	// 阻塞主 goroutine，直到 WaitGroup 的计数器归零
	// 也就是直到 4 个 sing() 全部执行完毕
	wait.Wait()

	// 所有子任务完成后，才会执行到这里
	fmt.Println("主线程结束")
}
