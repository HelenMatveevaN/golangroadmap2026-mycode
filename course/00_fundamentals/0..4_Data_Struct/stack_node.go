// СТЭК (Stack) — Принцип LIFO
package main

type Node struct {
	value int
	next  *Node
}

type ListStack struct {
	head *Node
}

// Push добавляет элемент в начало списка
func (s *ListStack) Push(val int) {
	newNode := &Node{value: val, next: s.head}
	s.head = newNode
}

// Pop удаляет первый элемент из списка
func (s *ListStack) Pop() (int, bool) {
	if s.head == nil {
		return 0, false
	}
	val := s.head.value
	s.head = s.head.next
	return val, true
}
