package main

import (
	"fmt"
	"os"
	"strconv"
)

/*
➜  0.6_first_prog git:(main) ✗ go run main.go 12 + 4
*/

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Ошибка: недостаточно аргументов")
		fmt.Println("Использование: go run main.go <число1> <операция> <число2>")
		fmt.Println("Пример: go run main.go 10 + 5")
		return
	}

	num1, err1 := strconv.ParseFloat(os.Args[1], 64)
	op := os.Args[2]
	num2, err2 := strconv.ParseFloat(os.Args[3], 64)

	if err1 != nil || err2 != nil {
		fmt.Println("Ошибка: введите корректные числа!")
		return
	}

	var result float64

	switch op {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		if num2 == 0 {
			fmt.Println("Ошибка: деление на ноль!")
			return
		}
		result = num1 / num2
	default:
			fmt.Println("Ошибка: неизвестная операция %s", op)
	}

	fmt.Printf("Результат: %g\n", result)
}