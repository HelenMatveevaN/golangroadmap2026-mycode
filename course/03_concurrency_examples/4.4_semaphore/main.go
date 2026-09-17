package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const maxConcurrent = 3 // не более 3 одновременных операций
	const totalJobs = 10

	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i := 1; i <= totalJobs; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem <- struct{}{} // захватываем семафор (блокирует если занято)
			defer func() { <-sem }() // освобождаем при выходе

			fmt.Printf("Задача %d начата (активных: %d)\n", id, len(sem))
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Задача %d завершена\n", id)
		}(i)
	}

	wg.Wait()
}