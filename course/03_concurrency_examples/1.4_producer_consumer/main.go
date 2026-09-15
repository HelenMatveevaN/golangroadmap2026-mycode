package main

import "fmt"

// generate возвращает канал, в который пишет числа в отдельной горутине
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// square читает из in и пишет квадраты в новый канал
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func main() {
	// Строим конвейер: generate → square → print
	gen := generate(2, 3, 4, 5)
	sq := square(gen)

	for result := range sq {
		fmt.Println(result) // 4, 9, 16, 25
	}
}