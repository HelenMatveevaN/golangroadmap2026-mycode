package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ПЛОХО: гонка данных — запустите с -race чтобы увидеть
func badCounter() int {
	counter := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(m *sync.Mutex) {
			defer wg.Done()
			mu.Lock()
			counter++ // DATA RACE: несинхронизированный доступ
			mu.Unlock()
		}(&mu)
	}

	wg.Wait()
	return counter
}

// ХОРОШО: через atomic
func goodCounterAtomic() int64 {
	var counter atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()
	return counter.Load()
}

// ХОРОШО: через канал
func goodCounterChannel() int {
	ch := make(chan struct{}, 1000)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- struct{}{}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	total := 0
	for range ch {
		total++
	}
	return total
}

func main() {
	fmt.Println("Atomic:", goodCounterAtomic())   // всегда 1000
	fmt.Println("Channel:", goodCounterChannel()) // всегда 1000
	fmt.Println("Bad:", badCounter()) // раскомментируй и запусти с -race
}