package main

import (
	"fmt"
	"sync"

	//"github.com/ChrisGora/semaphore"
	"math/rand"
	"time"
)

type buffer struct {
	b                 []int
	size, read, write int
	count             int //number of item currently in the buffer

	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
}

// This is a constructor-like function
// It creates a slice b of length size.
// It sets read and write to 0.
// It returns a buffer value (not a pointer).
func newBuffer(size int) *buffer {
	buf := &buffer{
		b:    make([]int, size),
		size: size,
	}
	buf.notFull = sync.NewCond(&buf.mu)
	buf.notEmpty = sync.NewCond(&buf.mu)
	return buf

}

func (buf *buffer) String() string {
	return fmt.Sprintf("&{%v %d %d %d}", buf.b, buf.size, buf.read, buf.write)
}

func (buf *buffer) get() int {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	//wait until there is something to read
	for buf.count == 0 {
		buf.notEmpty.Wait()
	}
	x := buffer.b[buffer.read]
	fmt.Println("Get\t", x, "\t", buffer)

	buf.read = (buf.read + 1) % buf.size
	buf.count--

	// One slot has become free.
	buf.notFull.Signal()

	return x
}

func (buf *buffer) put(x int) {
	buf.mu.Lock()
	defer buf.mu.Unlock()

	// Wait until there is a free slot.
	for buf.count == buf.size {
		buf.notFull.Wait()
	}

	buf.b[buf.write] = x
	fmt.Println("Put\t", x, "\t", buf)

	buf.write = (buf.write + 1) % buf.size
	buf.count++

	// One item is now available.
	buf.notEmpty.Signal()
}

func producer(buffer *buffer, start, delta int) {
	x := start
	for {
		buffer.put(x)
		x = x + delta
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func consumer(buffer *buffer) {
	for {
		_ = buffer.get()
		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	}
}

func main() {
	buf := newBuffer(5)

	go producer(buf, 1, 1)
	go producer(buf, 1000, -1)

	consumer(buf)
}
