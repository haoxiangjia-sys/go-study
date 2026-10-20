package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"strings"
)

// Message 表示一条来自某客户端的消息
type Message struct {
	sender  int
	message string
}

// handleError 统一处理非致命错误：只打印，不退出
func handleError(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}

// acceptConns 持续接受新连接，通过 channel 交给 main
func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
			continue
		}
		conns <- conn
	}
}

// handleClient 为单个客户端循环读取消息
func handleClient(client net.Conn, clientid int, msgs chan Message) {
	reader := bufio.NewReader(client)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("client %d disconnected: %v\n", clientid, err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		msgs <- Message{sender: clientid, message: line}
	}
}

func main() {
	// 1. 解析端口参数
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	// 2. 创建监听器
	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		log.Fatal(err) // 致命错误，退出
	}
	log.Println("Chat server listening on", *portPtr)

	// 3. 创建 channel 和共享状态
	conns := make(chan net.Conn)      // accept → main
	msgs := make(chan Message)        // handleClient → main
	clients := make(map[int]net.Conn) // ID → 连接
	nextID := 0

	// 4. 启动 accept goroutine
	go acceptConns(ln, conns)

	// 5. 主事件循环
	for {
		select {
		case conn := <-conns:
			// (a) 处理新客户端连接
			nextID++
			id := nextID
			clients[id] = conn
			log.Printf("client %d connected from %s\n", id, conn.RemoteAddr())
			go handleClient(conn, id, msgs)

		case msg := <-msgs:
			// (b) 处理客户端消息
			log.Printf("client %d says: %s\n", msg.sender, msg.message)
			// (c) 广播给除发送者外的所有客户端
			for id, conn := range clients {
				if id == msg.sender {
					continue
				}
				_, err := conn.Write([]byte(msg.message + "\n"))
				if err != nil {
					handleError(err)
					delete(clients, id)
				}
			}
		}
	}
}
