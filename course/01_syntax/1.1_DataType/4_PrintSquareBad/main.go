package main

import "fmt"

func PrintSquareBad(val any) int {
	v := val.(int)
	return v*v
}

func PrintSquareGood(val int) int {
	return val*val
}

func main() {
	v:= PrintSquareGood(5)
	fmt.Println("PrintSquareGood: ", v)

	v = PrintSquareBad("пять")
}