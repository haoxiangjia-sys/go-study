package main

import (
	"fmt"
	"sync"

	//"github.com/ChrisGora/semaphore"
	"math/rand"
	"time"

	"github.com/ChrisGora/semaphore"
)

type buffer struct {
	b                 []int
	size, read, write int
}

func newBuffer(size int) buffer {
	return buffer{
		b:     make([]int, size),
		size:  size,
		read:  0,
		write: 0,
	}
}

func (buffer *buffer) put(x int) {
	buffer.b[buffer.write] = x
	buffer.write = (buffer.write + 1) % buffer.size
}

func (buffer *buffer) get() int {
	x := buffer.b[buffer.read]
	buffer.read = (buffer.read + 1) % buffer.size
	return x
}

func producer(buffer *buffer, spaceAvailable, workAvailable semaphore.Semaphore, mutex *sync.Mutex, start, delta int) {
	x := start
	for {
		spaceAvailable.Wait()   // 1. 拿一个"空位"令牌，满了就等
		mutex.Lock()            // 2. 上锁
		buffer.put(x)           // 3. 写数据
		mutex.Unlock()          // 4. 解锁
		fmt.Println("Put\t", x) // 锁外打印
		workAvailable.Post()    // 5. 放一个"有货"令牌，通知消费者

		x = x + delta
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func consumer(buffer *buffer, spaceAvailable, workAvailable semaphore.Semaphore, mutex *sync.Mutex) {
	for {
		workAvailable.Wait() // 1. 拿一个"有货"令牌，空了就等
		mutex.Lock()
		_ = buffer.get()
		mutex.Unlock()
		spaceAvailable.Post() // 5. 放一个"空位"令牌，通知生产者

		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	}
}

func main() {
	buffer := newBuffer(5)
	var mutex sync.Mutex

	spaceAvailable := semaphore.Init(5, 5) // 一开始 5 个空位
	workAvailable := semaphore.Init(5, 0)  // 一开始 0 个数据

	go producer(&buffer, spaceAvailable, workAvailable, &mutex, 1, 1)
	go producer(&buffer, spaceAvailable, workAvailable, &mutex, 1000, -1)

	consumer(&buffer, spaceAvailable, workAvailable, &mutex)
}
