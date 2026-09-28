//package main
//
//import (
//	"fmt"
//	"sync"
//)
//
//var num int
//var wait sync.WaitGroup
//
//func add() {
//	for i := 0; i < 1000000; i++ {
//		num++
//	}
//	wait.Done()
//}
//func reduce() {
//	for i := 0; i < 1000000; i++ {
//		num--
//	}
//	wait.Done()
//}
//
//func main() {
//	wait.Add(2)
//	go add()
//	go reduce()
//	wait.Wait()
//	fmt.Println(num)
//
//} //因为这里存在数据竞争（data race），所以结果不可预测
//
//我们不能在并发模式下读写map
//
//如果要这样做
//
//给读写操作加锁
//使用sync.Map
//加锁

//package main
//
//import (
//	"fmt"
//	"sync"
//	"time"
//)
//
//var wait sync.WaitGroup
//var mp = map[string]string{}
//var lock sync.Mutex
//
//func reader() {
//	for {
//		lock.Lock()
//		fmt.Println(mp["time"])
//		lock.Unlock()
//	}
//	wait.Done()
//}
//func writer() {
//	for {
//		lock.Lock()
//		mp["time"] = time.Now().Format("15:04:05")
//		lock.Unlock()
//	}
//	wait.Done()
//}
//
//func main() {
//	wait.Add(2)
//	go reader()
//	go writer()
//	wait.Wait()
//}

package main

import (
	"fmt"
	"sync"
	"time"
)

var wait sync.WaitGroup
var mp = sync.Map{} //它是一个“并发安全的 map”

func reader() {
	for {

		fmt.Println(mp.Load("time"))
	}
	wait.Done()
}
func writer() {
	for {
		mp.Store("time", time.Now().Format("15:04:05"))
	}
	wait.Done()
}

func main() {
	wait.Add(2)
	go reader()
	go writer()
	wait.Wait()

}
