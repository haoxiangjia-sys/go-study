package main

import ("flag"
	"fmt"
	"net"
	"uk.ac.bris.cs/distributed3/pairbroker/stubs"
	"net/rpc")

var multiplyResults = make(chan int, 10) // Multiply 的结果送到这里
var client *rpc.Client                   // 到 broker 的连接，pipeline 也要用
type Factory struct {}

//(f *Factory)这个函数属于 Factory 类型 
//req stubs.Pair	第一个参数：broker 发来的任务，里面有 X 和 Y 两个数
//res *stubs.JobReport	第二个参数：指针，用来把结果交回给 broker
//(err error)	返回值：出没出错，不写内容默认是"没出错"
func (f *Factory) Multiply(req stubs.Pair, res *stubs.JobReport) (err error){
	res.Result = req.X * req.Y  //算乘法，结果放进 res 的 Result 里
	fmt.Println(req.X, "*", req.Y, "=", res.Result)
		multiplyResults <- res.Result
	return
}
func (f *Factory) Divide(req stubs.Pair, res *stubs.JobReport) (err error){
	if req.Y == 0 {
		fmt.Println("Cannot divide", req.X, "by 0, skipping")
		return
	}
	res.Result = req.X / req.Y  //算chu法，结果放进 res 的 Result 里
	fmt.Println(req.X, "/", req.Y, "=", res.Result)
	return
}


// 每收到两个乘法结果，就发布一道除法题
func pipeline() {
	for {
		x := <-multiplyResults   // 等第 1 个结果
		y := <-multiplyResults   // 等第 2 个结果
		req := stubs.PublishRequest{Topic: "divide", Pair: stubs.Pair{X: x, Y: y}}//填一张"发布单"：发到 divide 频道，内容是 x 和 y
		err := client.Call(stubs.Publish, req, new(stubs.StatusReport))//调用 broker 的 Publish 方法，就像 miner 做的那样
		if err != nil {
			fmt.Println("Publish error:", err)
		}
	}
}

func main(){
	pAddr := flag.String("ip", "127.0.0.1:8050", "IP and port to listen on")
	brokerAddr := flag.String("broker","127.0.0.1:8030", "Address of broker instance")
	flag.Parse()

	//1Factory 的方法可以被调用"。&Factory{} 是创建一个 Factory 并取它的地址
	rpc.Register(&Factory{})
	//在 127.0.0.1:8050 上"开门"，等别人连进来。返回一个 listener（门）和 err（有没有出错）
	listener, err := net.Listen("tcp", *pAddr)
	if err != nil {
		panic(err)
	}
	//等 main 结束时自动关门
	defer listener.Close()


	//2"拨号"连接 broker（127.0.0.1:8030），返回一条连接 client
	client, err = rpc.Dial("tcp", *brokerAddr)
	if err != nil {
		panic(err)
	}
	defer client.Close()


	// 确保 divide 频道存在（这次没有 divide miner 来创建它了）
	client.Call(stubs.CreateChannel, stubs.ChannelRequest{Topic: "divide", Buffer: 10}, new(stubs.StatusReport))

	// 3. 订阅 multiply 频道
	subscription := stubs.Subscription{
		Topic:          "multiply",
		FactoryAddress: *pAddr,
		Callback:       "Factory.Multiply",
	}
	status := new(stubs.StatusReport)                         // 准备一个地方放 broker 的回复
    err = client.Call(stubs.Subscribe, subscription, status)  // 调用 broker 的 Subscribe 方法
	if err != nil {
		panic(err)
	}
	fmt.Println("Subscribed to multiply, listening on", *pAddr)

	// 再订阅 divide 频道
	divideSub := stubs.Subscription{
		Topic:          "divide",
		FactoryAddress: *pAddr,
		Callback:       "Factory.Divide",
	}
	err = client.Call(stubs.Subscribe, divideSub, status)
	if err != nil {
		panic(err)
	}
	fmt.Println("Subscribed to divide")


	// 4. 启动流水线，然后一直接收 broker 派来的活
	go pipeline()
	rpc.Accept(listener)
}


