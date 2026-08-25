package main

import "fmt"

func testDefer() {
	x := 10
	defer fmt.Println("Значение в defer:", x)
	x = 20
	fmt.Println("Значение в конце функции:", x)
}

func main() {
	testDefer()
}