package main

import "fmt"

func sum(numbers []int, result chan<- int) {
	total := 0
	for _, n := range numbers {
		total += n
	}
	result <- total // отправка блокирует до момента получения
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ch := make(chan int) // небуферизованный канал

	// Считаем сумму первой и второй половины параллельно
	go sum(numbers[:5], ch)
	go sum(numbers[5:], ch)

	// Получаем два результата
	a, b := <-ch, <-ch
	fmt.Println("Сумма:", a+b) // 55
}