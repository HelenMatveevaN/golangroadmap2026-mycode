package main

import (
	"fmt"
	"sync"
)

// fanOut распределяет задачи по N горутинам
func fanOut(input <-chan int, n int) []<-chan int {
	outputs := make([]<-chan int, n)
	for i := 0; i < n; i++ {
		out := make(chan int)
		outputs[i] = out
		go func(ch chan<- int) {
			defer close(ch)
			for v := range input {
				ch <- v * v // каждая горутина возводит в квадрат
			}
		}(out)
	}
	return outputs
}

// fanIn сливает несколько каналов в один
func fanIn(channels ...<-chan int) <-chan int {
	merged := make(chan int)
	var wg sync.WaitGroup

	output := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			merged <- v
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 8; i++ {
			input <- i
		}
	}()

	// Раздаём 3 воркерам и собираем обратно
	workers := fanOut(input, 3)
	results := fanIn(workers...)

	for v := range results {
		fmt.Println(v)
	}
}