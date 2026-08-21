// ОЧЕРЕДЬ (Queue) — Принцип FIFO (на связном списке)
package main

type ListQueue struct {
	head *Node
	tail *Node
}

// Enqueue добавляет элемент в хвост (tail)
func (q *ListQueue) Enqueue(val int) {
	newNode := &Node{value: val, next: nil}
	if q.tail == nil { //проверка на пустую очередь
		q.head = newNode //элемент становится началом и концом очереди одновременно
		q.tail = newNode
		return
	}
	q.tail.next = newNode //если в очереди уже есть элементы
	q.tail = newNode //переопределяем хвост
}

// Dequeue забирает элемент из головы (head)
func (q *ListQueue) Dequeue() (int, bool) {
	if q.head == nil {
		return 0, false
	}
	val := q.head.value
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil //если очередь опустела, зануляем хвост
	}
	return val, true
}