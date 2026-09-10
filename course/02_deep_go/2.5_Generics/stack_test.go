package main

import "testing"

/*
go test -bench=. -benchmem
*/

// Бенчмарк для Дженерик-стека
func BenchmarkGenericStack(b *testing.B) {
	s := &GenericStack[int]{}
	b.ResetTimer() // Сбрасываем таймер перед тестом
	
	for i := 0; i < b.N; i++ {
		s.Push(i)
		_ = s.Pop()
	}
}

// Бенчмарк для старого стека на any (интерфейсах)
func BenchmarkInterfaceStack(b *testing.B) {
	s := &InterfaceStack{}
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		s.Push(i) // Здесь происходит боксинг int -> any (выделение памяти в куче)
		val := s.Pop()
		_ = val.(int) // Вынужденное утверждение типа (Type Assertion)
	}
}
