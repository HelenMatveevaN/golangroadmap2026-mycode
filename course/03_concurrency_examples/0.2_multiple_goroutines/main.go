package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1) // сообщаем WaitGroup, что появилась ещё одна горутина
		go func(n int) {
			defer wg.Done() // сигнализируем о завершении
			fmt.Printf("Горутина %d работает\n", n)
		}(i)
	}

	wg.Wait() // ждём пока все горутины вызовут Done()
	fmt.Println("Все горутины завершились")
}