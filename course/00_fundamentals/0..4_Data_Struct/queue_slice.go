// ОЧЕРЕДЬ (Queue) — Принцип FIFO (на слайсе)
package main

type ArrayQueue struct {
	data     []int
	head     int
	tail     int
	size     int
	capacity int
}

func NewArrayQueue(capacity int) *ArrayQueue {
	return &ArrayQueue{
		data:     make([]int, capacity),
		capacity: capacity,
	}
}

// Enqueue добавляет элемент в конец очереди
func (q *ArrayQueue) Enqueue(val int) bool {
	if q.size == q.capacity {
		return false //очередь переполнена
	}
	q.data[q.tail] = val
	q.tail = (q.tail + 1) % q.capacity //циклич-ий сдвиг хвоста
	q.size++
	return true
}

func (q *ArrayQueue) Dequeue() (int, bool) {
	if q.size == 0 {
		return 0, false //очередь пуста
	}
	val := q.data[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--
	return val, true
}
