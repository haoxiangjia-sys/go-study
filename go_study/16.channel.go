//var c chan int // 声明一个传递 int 的通道，此时是 nil
//
//c = make(chan int)      // 无缓冲通道
//c = make(chan int, 1)   // 有缓冲通道，缓冲大小为 1
//
//c <- 1        // 发送：把 1 放进通道
//n := <-c      // 接收：从通道取出一个值
//n, ok := <-c  // 接收并判断通道是否已关闭
//
//close(c)

package main

import "fmt"

func main() {
	var c chan int // 声明，此时 c == nil

	c = make(chan int, 1) // 初始化，缓冲大小 1
	c <- 1                // 缓冲区为空，可以放进去，不阻塞

	// c <- 2 // 缓冲区满了，没人接收，会死锁
	fmt.Println(<-c) // 取出 1，缓冲区又空了

	// fmt.Println(<-c) // 缓冲区为空，没人发送，会死锁

	c <- 2
	n, ok := <-c       //一个 bool，表示这次接收是否成功。
	fmt.Println(n, ok) // 输出 2 true

	close(c) // 关闭通道

	// c <- 3 // panic: send on closed channel
	// n, ok = <-c // 关闭后可以读，返回 0 false
	fmt.Println(c)
}



//1，数据传递
ch := make(chan int)
go func() {ch <- 42}()
fmt.Println(<-ch)

//2. 同步信号
//channel 也可以只用来“通知”，不传具体数据：

go
done := make(chan struct{})
go func() {
	// 干活
	close(done) // 或者 done <- struct{}{}
}()
<-done
struct{} 是空结构体，不占内存，常用于纯信号。

//3. 遍历

for v := range ch {
fmt.Println(v)
}
//4. 超时控制
go
select {
case v := <-ch:
fmt.Println("收到", v)
case <-time.After(2 * time.Second):
fmt.Println("超时")
}
//5. 多路复用（select）
go
select {
case v := <-ch1:
fmt.Println("来自 ch1", v)
case v := <-ch2:
fmt.Println("来自 ch2", v)
}

























