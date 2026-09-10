package main

// Дженерик-стек (Новый подход)
type GenericStack[T any] struct {
	elements []T
}

func (s *GenericStack[T]) Push(val T) {
	s.elements = append(s.elements, val)
}

func (s *GenericStack[T]) Pop() T {
	if len(s.elements) == 0 {
		var zero T
		return zero
	}
	index := len(s.elements) - 1
	val := s.elements[index]
	s.elements = s.elements[:index]
	return val
}

// Стек на интерфейсах any (Старый подход)
type InterfaceStack struct {
	elements []any
}

func (s *InterfaceStack) Push(val any) {
	s.elements = append(s.elements, val)
}

func (s *InterfaceStack) Pop() any {
	if len(s.elements) == 0 {
		return nil
	}
	index := len(s.elements) - 1
	val := s.elements[index]
	s.elements = s.elements[:index]
	return val
}
