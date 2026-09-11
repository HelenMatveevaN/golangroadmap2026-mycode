package containers

import "testing"

//go test -v -cover

func TestStack(t *testing.T) {
	stack := NewStack[int]()

	if !stack.IsEmpty() {
		t.Error("Ожидался пустой стек при инициализации")
	}

	stack.Push(10)
	stack.Push(20)

	if stack.IsEmpty() {
		t.Error("Стек не должен быть пустым после Push")
	}

	// Тестируем Pop (LIFO order)
	val, ok := stack.Pop()
	if !ok || val != 20 {
		t.Errorf("Ожидалось извлечение 20, получено %v (ok: %v)", val, ok)
	}

	val, ok = stack.Pop()
	if !ok || val != 10 {
		t.Errorf("Ожидалось извлечение 10, получено %v (ok: %v)", val, ok)
	}

	_, ok = stack.Pop()
	if ok {
		t.Error("Ожидалось ok=false при извлечении из пустого стека")
	}

	if !stack.IsEmpty() {
		t.Error("Стек должен снова стать пустым")
	}
}

func TestQueue(t *testing.T) {
	queue := NewQueue[string]()

	if !queue.IsEmpty() {
		t.Error("Ожидалась пустая очередь при инициализации")
	}

	queue.Enqueue("first")
	queue.Enqueue("second")

	if queue.IsEmpty() {
		t.Error("Очередь не должна быть пустой после Enqueue")
	}

	// Тестируем Dequeue (FIFO order)
	val, ok := queue.Dequeue()
	if !ok || val != "first" {
		t.Errorf("Ожидалось извлечение 'first', получено %v (ok: %v)", val, ok)
	}

	val, ok = queue.Dequeue()
	if !ok || val != "second" {
		t.Errorf("Ожидалось извлечение 'second', получено %v (ok: %v)", val, ok)
	}

	if !queue.IsEmpty() {
		t.Error("Очередь должна быть пустой")
	}
}

func TestSet(t *testing.T) {
	set := NewSet[int]()

	set.Add(5)
	set.Add(5)
	set.Add(10)

	if set.Size() != 2 {
		t.Errorf("Ожидался размер 2, получили %d", set.Size())
	}

	if !set.Has(5) {
		t.Error("Элемент 5 должен быть в множестве")
	}

	set.Remove(5)
	if set.Has(5) {
		t.Error("Элемент 5 не должен быть в множестве после удаления")
	}
}