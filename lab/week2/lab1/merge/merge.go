package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"sync"
)

// merge takes two sorted sub-arrays from slice and sorts them.
// The resulting array is put back in slice.
func merge(slice []int32, middle int) {
	sliceClone := make([]int32, len(slice))
	copy(sliceClone, slice)
	a := sliceClone[middle:]
	b := sliceClone[:middle]
	i := 0
	j := 0
	for k := 0; k < len(slice); k++ {
		if i >= len(a) {
			slice[k] = b[j]
			j++
		} else if j >= len(b) {
			slice[k] = a[i]
			i++
		} else if a[i] > b[j] {
			slice[k] = b[j]
			j++
		} else {
			slice[k] = a[i]
			i++
		}
	}
}

// Sequential merge sort.
func mergeSort(slice []int32) {
	if len(slice) > 1 {
		middle := len(slice) / 2
		mergeSort(slice[:middle])
		mergeSort(slice[middle:])
		merge(slice, middle)
	}
}

//// TODO: Parallel merge sort.
//func parallelMergeSort(slice []int32) {
//	//1如果 len(slice) <= 1，直接返回。
//	if len(slice) <= 1 {
//		return
//	}
//	//2找到中点 middle。
//	middle := len(slice) / 2
//	//3启动两个 goroutine：
//	var wg sync.WaitGroup
//	wg.Add(2)
//	//4一个排序左半段 slice[:middle]
//	go func() {
//		defer wg.Done()
//		parallelMergeSort(slice[:middle])
//	}()
//	//5一个排序右半段 slice[middle:]
//	go func() {
//		defer wg.Done()
//		parallelMergeSort(slice[middle:])
//	}()
//	//6等待两个 goroutine 都完成。
//	wg.Wait()
//	//7调用 merge(slice, middle) 合并。
//	merge(slice, middle)
//}
//改进：每次分割只创建一个新 goroutine
//让当前 goroutine 继续处理右半区，只把左半区交给新 goroutine：

func parallelMergeSort(slice []int32) {
	if len(slice) <= 1 {
		return
	}
	middle := len(slice) / 2

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		parallelMergeSort(slice[:middle]) // 新 goroutine 处理左半
	}()

	parallelMergeSort(slice[middle:]) // 当前 goroutine 处理右半
	wg.Wait()
	merge(slice, middle)
}

// main starts tracing and in parallel sorts a small slice.
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

	slice := make([]int32, 0, 100)
	for i := int32(100); i > 0; i-- {
		slice = append(slice, i)
	}

	parallelMergeSort(slice)
	fmt.Println(slice)
}
