package main

import "fmt"

type User struct {
	Name string
}

//go build -gcflags="-m" main.go

//Inlining (встраивание) 

// Сценарий 1: Возвращаем значение (Value Semantics)
func CreateUserByValue() User {
	u := User{Name: "Alice"}
	return u
}

// Сценарий 2: Возвращаем указатель (Pointer Semantics)
func CreateUserByPointer() *User {
	u := User{Name: "Bob"}
	return &u
}

func main() {
	user1 := CreateUserByValue()
	_ = user1

	user2 := CreateUserByPointer()
	_ = user2

	// Сценарий 3: Передача в прожорливую встроенную функцию
	x := 42
	fmt.Println(x)
}