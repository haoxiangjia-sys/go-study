//1
//package main
//
//import (
//	"fmt"
//	"sync"
//)
//
//func main() {
//	sum := 0
//	var wg sync.WaitGroup //counter
//	var mu sync.Mutex
//
//	for i := 0; i < 1000; i++ {
//		wg.Add(1) //计数器 +1，表示"又多了一个任务要等"
//		go func() {
//			mu.Lock() //// 拿到锁才能进入
//			sum = sum + 1
//			mu.Unlock() //用完放开，下一个才能进
//			wg.Done() //计数器 -1，表示"又多了一个任务要等"
//		}()
//	}
//
//	wg.Wait()
//	fmt.Println(sum)
//}

//2
//package main
//
//import "fmt"
//
//func main() {
//	ch := make(chan int)
//
//	for i := 0; i < 1000; i++ {
//		go func() {
//			ch <- 1
//		}()
//	}
//
//	sum := 0
//	for i := 0; i < 1000; i++ {
//		sum += <-ch
//	}
//
//	fmt.Println(sum)
//}

// 3
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var sum int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			atomic.AddInt64(&sum, 1)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(sum)
}
