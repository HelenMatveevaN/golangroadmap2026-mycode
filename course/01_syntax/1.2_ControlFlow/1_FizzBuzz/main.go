package main

import "fmt"

func main() {
	// Итерируемся по числам от 1 до 20 включительно
	for i := 1; i <= 20; i++ {

		// Усложнение: if с инициализацией для проверки на четность
		if remainder := i % 2; remainder == 0 {
			fmt.Printf("[%d — Четное] ", i)
		} else {
			fmt.Printf("[%d — Нечетное] ", i)
		}
		//fmt.Printf(remainder) - ./main.go:16:14: undefined: remainder

		// Усложнение: switch без выражения (работает как аналог if-else if)
		switch {
		case i%3 == 0 && i%5 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}
}