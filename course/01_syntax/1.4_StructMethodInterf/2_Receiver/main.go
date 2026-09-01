package main

import "fmt"

//Разница ресиверов: Написать два метода для структуры Counter (счетчик): 
//один с value-ресивером, другой с pointer-ресивером, 
//и увидеть, почему один из них не увеличивает внутренний счетчик.

type Counter struct {
	value int
}

// 1. Метод с Value-ресивером (передача по значению)
func (c Counter) IncrementByValue() {
	c.value++
	fmt.Printf("[Внутри Value-метода] c.value = %d (адрес в памяти: %p)\n", c.value, &c)
}

// 2. Метод с Pointer-ресивером (передача по указателю)
func (c *Counter) IncrementByPointer() {
	c.value++
	fmt.Printf("[Внутри Pointer-метода] c.value = %d (адрес в памяти: %p)\n", c.value, c)
}

func main() {
	myCounter := Counter{value: 0}
	fmt.Printf("[В main — Старт] Адрес оригинального счетчика в памяти: %p\n\n", &myCounter)

	// --- ТЕСТ 1: Вызываем Value-ресивер ---
	fmt.Println("=== Вызываем IncrementByValue ===")
	myCounter.IncrementByValue()
	fmt.Println("[В main — После вызова] Реальное значение счетчика myCounter.value:", myCounter.value) 
	fmt.Println()

	// --- ТЕСТ 2: Вызываем Pointer-ресивер ---
	fmt.Println("=== Вызываем IncrementByPointer ===")
	myCounter.IncrementByPointer()
	fmt.Println("[В main — После вызова] Реальное значение счетчика myCounter.value:", myCounter.value) 
	fmt.Println()
}