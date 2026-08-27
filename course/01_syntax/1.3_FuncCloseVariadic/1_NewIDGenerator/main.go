package main

import "fmt"

func NewIDGenerator() func() int {
	id := 0

	return func() int {
		id++
		return id
	}
}

func main() {
	generateUserID := NewIDGenerator()

	fmt.Println("Генератор пользователей:")
	fmt.Println(generateUserID())
	fmt.Println(generateUserID())
	fmt.Println(generateUserID())

	generateOrderID := NewIDGenerator()
	fmt.Println("\nГенератор заказов:")
	fmt.Println(generateOrderID())
	fmt.Println(generateOrderID())	
}