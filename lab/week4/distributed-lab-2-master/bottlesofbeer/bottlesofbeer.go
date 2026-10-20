package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"
)

var nextAddr string

// rpc.Client 是"一条已经接通的电话线"
// 前面的 * 表示指针，可以是 nil，也就是"还没接通"
var nextClient *rpc.Client

// Request carries how many bottles are left when it's this buddy's turn.
type Request struct {
	Bottles int
}

type Response struct{}

// Buddy is the RPC type every process registers. All buddies run the same code.
type Buddy struct{}

// Sing is called by the previous buddy: sing one verse, then pass the turn on.
func (b *Buddy) Sing(req Request, res *Response) (err error) {
	n := req.Bottles
	if n == 0 {
		fmt.Println("No more bottles of beer on the wall, no more bottles of beer.")
		fmt.Println("Go to the store and buy some more... 🍺 Song finished!")
		return
	}

	fmt.Printf("%d %s of beer on the wall, %d %s of beer.\n", n, bottleWord(n), n, bottleWord(n))
	fmt.Printf("Take one down, pass it around, %d %s of beer on the wall.\n\n", n-1, bottleWord(n-1))
	time.Sleep(500 * time.Millisecond) // so you can watch the turns move round

	passOn(n - 1)
	return
}

func bottleWord(n int) string {
	if n == 1 {
		return "bottle"
	}
	return "bottles"
}

// passOn asks the next buddy to sing, WITHOUT waiting for the reply (client.Go),
// so calls don't pile up in a chain all the way round the ring.
func passOn(n int) {
	if nextClient == nil {
		c, err := rpc.Dial("tcp", nextAddr)
		if err != nil {
			panic(err)
		}
		nextClient = c
	}
	nextClient.Go("Buddy.Sing", Request{Bottles: n}, new(Response), nil)
}

func main() {
	// ① 读命令行参数

	//flag.String 的三个参数：参数名、默认值、说明。它返回的是一个指针，所以后面用的时候要写 *thisPort 才能拿到字符串。
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	//tringVar 是另一种写法：不返回指针，而是直接把结果写进你给的变量。&nextAddr 是全局变量的地址。用这种写法是因为 nextAddr 要在 passOn 里用，必须是全局的。
	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for next member of the round.")

	//flag.String 的三个参数：参数名、默认值、说明
	bottles := flag.Int("n", 0, "Bottles of Beer (launches song if not 0)")
	//flag.Parse() 必须调用，它才会真正去读命令行
	flag.Parse()

	// ② 开门营业（当服务器）
	rpc.Register(&Buddy{})
	listener, err := net.Listen("tcp", ":"+*thisPort)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Listening on port", *thisPort, "- next buddy is", nextAddr)

	// ③ 如果我是领唱，等回车后开始
	if *bottles > 0 {
		go func() {
			fmt.Println("Press Enter once every buddy is running to start the song...")
			bufio.NewReader(os.Stdin).ReadString('\n')
			new(Buddy).Sing(Request{Bottles: *bottles}, new(Response))
		}()
	}

	// ④ 一直等别人打电话进来
	rpc.Accept(listener)
}
