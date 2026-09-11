package containers
//LIFO (Last In, First Out)

// Stack — дженерик-структура. 
// Ключевое слово `any` означает, что вместо T может быть подставлен любой тип.
type Stack[T any] struct {
	elements []T
}

// NewStack — функция-конструктор для создания чистого стека.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		elements: make([]T, 0),
	}
}

// Push добавляет элемент на вершину стека.
// Обрати внимание: мы используем указатель (*Stack[T]), чтобы изменять внутреннее состояние структуры.
func (s *Stack[T]) Push(val T) {
	s.elements = append(s.elements, val)
}

// Pop удаляет и возвращает элемент с вершины стека.
// Если стек пуст, возвращает нулевое значение типа T и false.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}

	//last idx
	lastIdx := len(s.elements) - 1
	//запомним элемент
	val := s.elements[lastIdx]

	//слайс без последнего элемента
	s.elements = s.elements[:lastIdx]

	return val, true
}

// IsEmpty возвращает true, если стек пуст.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.elements)
}