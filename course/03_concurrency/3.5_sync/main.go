package main

//потокобезопасный счётчик через Mutex, 
//затем через atomic — сравни производительность.

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type MutexCounter struct {
	mu	  sync.Mutex
	value int64
}

type AtomicCounter struct {
	value int64
}

func (c *MutexCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func main() {
	myCounter := &MutexCounter{}

	var wg sync.WaitGroup

	for i:=0; i<1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			myCounter.Increment()
		}()
	}

	myCounter2 := &AtomicCounter{}
	var wg2 sync.WaitGroup

	for i:=0; i<1000; i++ {
		wg2.Add(1)

		go func() {
			defer wg2.Done()
			myCounter2.Increment()
		}()
	}

	wg.Wait()
	fmt.Println("MutexCounter.value = ", myCounter.value)

	wg2.Wait()
	fmt.Println("AtomicCounter.value = ", myCounter2.value)
}