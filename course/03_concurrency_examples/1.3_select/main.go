package main

import (
	"fmt"
	"time"
)

func fastWorker(ch chan<- string) {
	time.Sleep(100 * time.Millisecond)
	ch <- "быстрый результат"
}

func slowWorker(ch chan<- string) {
	time.Sleep(500 * time.Millisecond)
	ch <- "медленный результат"
}

func main() {
	fast := make(chan string, 1)
	slow := make(chan string, 1)

	go fastWorker(fast)
	go slowWorker(slow)

	// Ждём любого из двух или таймаут
	select {
	case result := <-fast:
		fmt.Println("Первым ответил:", result)
	case result := <-slow:
		fmt.Println("Первым ответил:", result)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("Таймаут!")
	}
}