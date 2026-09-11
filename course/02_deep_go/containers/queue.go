package containers

// Queue — дженерик-структура для очереди
type Queue[T any] struct {
	elements []T
}

// NewQueue — конструктор для создания очереди.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		elements: make([]T, 0),
	}
}

// Enqueue добавляет элемент в конец очереди
func (q *Queue[T]) Enqueue(val T) {
	q.elements = append(q.elements, val)
}

// Dequeue удаляет и возвращает первый элемент из начала очереди.
// Если очередь пуста, возвращает нулевое значение T и false.
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.elements) == 0 {
		var zero T
		return zero, false
	}

	//запомним элемент
	val := q.elements[0]

	//слайс без первого элемента
	q.elements = q.elements[1:]

	return val, true
}

// IsEmpty возвращает true, если очередь пуста
func (q *Queue[T]) IsEmpty() bool {
	return len(q.elements) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.elements)
}