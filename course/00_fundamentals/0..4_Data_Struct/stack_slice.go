// СТЭК (Stack) — Принцип LIFO
package main

type ArrayStack struct {
	data []int
}

// Push добавляет элемент на вершину стека
func (s *ArrayStack) Push(val int) {
	s.data = append(s.data, val)
}

// Pop удаляет и возвращает верхний элемент
func (s *ArrayStack) Pop() (int, bool) {
	if len(s.data) == 0 {
		return 0, false //стек пуст
	}
	lastIdx := len(s.data) - 1
	val := s.data[lastIdx]
	s.data = s.data[:lastIdx] //обрезаем последний элемент
	return val, true
}
