package main

import "fmt"

func main() {
	// Канал с буфером на 3 элемента
	ch := make(chan string, 3)

	ch <- "первый"  // не блокирует
	ch <- "второй"  // не блокирует
	ch <- "третий"  // не блокирует
	// ch <- "четвёртый" // заблокировало бы — буфер полон

	close(ch) // закрываем, чтобы range завершился

	for msg := range ch {
		fmt.Println(msg)
	}
}