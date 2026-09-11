package containers

// Set — дженерик-структура для множества уникальных элементов.
// Ограничение `comparable` разрешает использовать только те типы, 
// которые можно сравнивать через == (числа, строки, указатели).
type Set[T comparable] struct {
	elements map[T]struct{}
}

// NewQueue — конструктор для создания множества.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(map[T]struct{}),
	}
}

// Add добавляет элемент в множество. 
// Если он уже есть, map просто перезапишет ключ — уникальность сохраняется.
func (s *Set[T]) Add(val T) {
	s.elements[val] = struct{}{}
}

// Если он уже есть, map просто перезапишет ключ — уникальность сохраняется.
func (s *Set[T]) Has(val T) bool {
	_, exists := s.elements[val]
	return exists
}

// Remove удаляет элемент из множества
func (s *Set[T]) Remove(val T) {
	delete(s.elements, val)
}

// Size возвращает количество элементов в множестве
func (s *Set[T]) Size() int {
	return len(s.elements)
}