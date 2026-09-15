package main

import (
	"fmt"
	"time"
)

func main() {
	// Запускаем горутину — лёгкий поток Go-рантайма
	go func() {
		fmt.Println("Привет из горутины!")
	}()

	// Без этого main завершится раньше, чем горутина успеет напечатать
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Привет из main!")
}