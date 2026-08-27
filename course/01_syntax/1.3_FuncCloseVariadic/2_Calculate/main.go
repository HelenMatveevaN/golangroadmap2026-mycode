package main

import "fmt"

func Calculate(a, b int, op func(int, int) int) int {
	return op(a, b)
}

func main() {
	x, y := 10, 5

	sumResult := Calculate(x, y, func(num1, num2 int) int {
		return num1 + num2
	})
	fmt.Printf("Сложение: %d + %d = %d\n", x, y, sumResult)

	mulResult := Calculate(x, y, func(num1, num2 int) int {
		return num1 * num2
	})
	fmt.Printf("Умножение: %d * %d = %d\n", x, y, mulResult)

	substract := func(num1, num2 int) int {
		return num1 - num2
	}
	subResult := Calculate(x, y, substract)
	fmt.Printf("Вычитание: %d - %d = %d\n", x, y, subResult)
}