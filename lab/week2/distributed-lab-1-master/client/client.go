package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

// read 持续从服务器读消息并显示
func read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nServer closed connection:", err)
			os.Exit(0)
		}
		fmt.Print(line) // 服务器消息自带 \n
	}
}

// write 持续读用户输入并发送
func write(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}

func main() {
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()

	// 连接服务器
	conn, err := net.Dial("tcp", *addrPtr)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()
	fmt.Println("Connected to", *addrPtr, "— type messages and press Enter")

	// 读放 goroutine，写放主线程
	go read(conn)
	write(conn)
}
