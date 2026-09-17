package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		requests <- i
	}
	close(requests)

	// Тикер — сигнализирует каждые 200ms (5 req/s)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for req := range requests {
		<-ticker.C
		fmt.Printf("Обрабатываем запрос %d в %s\n", req, time.Now().Format("15:04:05.000"))
	}
}