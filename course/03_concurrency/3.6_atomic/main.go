package main

import (
	"fmt"
	"time"
	"sync/atomic"
)

func main() {

	var IsRunning atomic.Bool
	IsRunning.Store(true) //running

	go func() {
		time.Sleep(500 * time.Millisecond)
		IsRunning.Store(false)
	}()

	ticker := time.Tick(200 * time.Millisecond)
	for range ticker {
		status := IsRunning.Load()
		fmt.Println(status)

		if status == false {
			break
		}
	}

	fmt.Println("[Main] Программа успешно завершена.")
}