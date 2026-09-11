package main

import "fmt"

/*
go build -gcflags="-m -l" main.go

Флаг -m включает вывод анализа памяти, 
флаг -l (lowercase L) запрещает компилятору делать инлайнинг (встраивание кода функций), 
чтобы мы увидели чистый escape-анализ для каждой функции.

*/

// 1. Остается на стеке: простое значение, копия уходит наружу
func stayOnStack() int {
	x := 42
	return x
}

// 2. Убегает: возвращаем указатель на локальную переменную
func escapePyPointer() *int {
	x := 100
	return &x
}

// 3. Убегает: замыкание хранит состояние переменной count
func escapeByClosure() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// 4. Убегает: fmt.Println принимает интерфейс any
func escapeToInterface() {
	y := 200
	fmt.Println(y)
}

// 5. Убегает: размер слайса динамический, компилятор перестраховывается
func escapeByDynamicSlice(n int) []int {
	s := make([]int, n)
	return s
}

func main() {
	_ = stayOnStack()
	_ = escapePyPointer

	f := escapeByClosure
	_ = f()

	escapeToInterface()
	_ = escapeByDynamicSlice(10)
}