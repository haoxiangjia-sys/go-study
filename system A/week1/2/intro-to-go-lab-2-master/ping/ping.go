package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

func foo(channel chan string) {
	for {
		// 一轮 ping-pong
		fmt.Println("Foo is sending: ping")
		channel <- "ping" // 发送 ping，阻塞直到 bar 接收
		pong := <-channel // 等待接收 pong
		fmt.Println("Foo has received:", pong)
	}
}

func bar(channel chan string) {
	for {
		// 一轮 ping-pong
		ping := <-channel // 等待接收 ping
		fmt.Println("Bar has received:", ping)
		fmt.Println("Bar is sending: pong")
		channel <- "pong" // 发送 pong，阻塞直到 foo 接收
	}
}

func pingPong() {
	channel := make(chan string) // 创建无缓冲字符串通道
	go foo(channel)              // 把 channel 传给 foo
	go bar(channel)              // 把 channel 传给 bar
	time.Sleep(500 * time.Millisecond)
}
func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	pingPong()
}
