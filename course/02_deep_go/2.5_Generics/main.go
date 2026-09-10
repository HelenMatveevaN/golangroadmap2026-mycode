package main

import (
	"errors"
	"fmt"
)

// ==========================================
// 1. Stack[T any] — Стек (LIFO: последним пришел — первым ушел)
// ==========================================

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Push(val T) {
	s.elements = append(s.elements, val)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.elements) == 0 {
		var zero T
		return zero, errors.New("stack is empty")
	}
	index := len(s.elements) - 1
	val := s.elements[index]
	s.elements = s.elements[:index]
	return val, nil
}

func (s *Stack[T]) Size() int {
	return len(s.elements)
}

// ==========================================
// 2. Queue[T any] — Очередь (FIFO: первым пришел — первым ушел)
// ==========================================
type Queue[T any] struct {
	elements []T
}

func (q *Queue[T]) Enqueue(val T) {
	q.elements = append(q.elements, val)
}

func (q *Queue[T]) Dequeue() (T, error) {
	if len(q.elements) == 0 {
		var zero T
		return zero, errors.New("queue is empty")
	}
	val := q.elements[0]
	q.elements = q.elements[1:]
	return val, nil
}

func (q *Queue[T]) Size() int {
	return len(q.elements)
}

// ==========================================
// 3. Set[T comparable] — Множество уникальных элементов
// Ограничение comparable обязательно, так как элементы будут ключами мапы.
// ==========================================
type Set[T comparable] struct {
	// Использование struct{} вместо bool экономит память (struct{} занимает 0 байт)
	data map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{data: make(map[T]struct{})}
}

func (s *Set[T]) Add(val T) {
	s.data[val] = struct{}{}
}

func (s *Set[T]) Has(val T) bool {
	_, exists := s.data[val]
	return exists
}

func (s *Set[T]) Delete(val T) {
	delete(s.data, val)
}

// ==========================================
// 4. CustomMap[K comparable, V any] — Обертка над встроенной мапой
// Демонстрирует использование двух разных параметров типов (K и V).
// ==========================================
type CustomMap [K comparable, V any] struct {
	store map[K]V
}

func NewCustomMap[K comparable, V any]() *CustomMap[K, V] {
	return &CustomMap[K, V]{store: make(map[K]V)}
}

func (m *CustomMap[K, V]) Set(key K, val V) {
	m.store[key] = val
}

func (m *CustomMap[K, V]) Get(key K) (V, bool) {
	val, exists := m.store[key]
	return val, exists
}

// ==========================================
// Тестирование контейнеров
// ==========================================
func main() {
	fmt.Println("--- 1. Тест Стека (Stack[int]) ---")
	intStack := &Stack[int]{}
	intStack.Push(10)
	intStack.Push(20)
	fmt.Println("Размер стека:", intStack.Size()) // 2
	val, _ := intStack.Pop()
	fmt.Println("Извлечено из стека (должно быть 20):", val)

	fmt.Println("\n--- 2. Тест Очереди (Queue[string]) ---")
	strQueue := &Queue[string]{}
	strQueue.Enqueue("Первый")
	strQueue.Enqueue("Второй")
	msg, _ := strQueue.Dequeue()
	fmt.Println("Извлечено из очереди (должно быть Первый):", msg)

	fmt.Println("\n--- 3. Тест Множества (Set[int]) ---")
	mySet := NewSet[int]()
	mySet.Add(5)
	mySet.Add(5)
	fmt.Println("Содержит 5?:", mySet.Has(5)) // true
	fmt.Println("Содержит 5?:", mySet.Has(10)) // true

	fmt.Println("\n--- 4. Тест Кастомной Мапы (CustomMap[string, float64]) ---")
	currencyMap := NewCustomMap[string, float64]()
	currencyMap.Set("USD/RUB", 91.50)
	rate, _ := currencyMap.Get("USD/RUB")
	fmt.Printf("Курс USD/RUB: %.2f\n", rate)
}