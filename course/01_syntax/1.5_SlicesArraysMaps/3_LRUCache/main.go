package main

import (
	"container/list"
	"fmt"
)

// pair будет храниться в поле Value внутри list.Element
type pair struct {
	key 	string
	value 	int
}

type LRUCache struct {
	capacity int 						//макс.размер кэша
	items	map[string]*list.Element	//быстр.поиск узла за О(1). Мапа указ-лей
	queue	*list.List					//2-связный список
}

// NewLRUCache — конструктор кэша
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity:	capacity,
		items:		make(map[string]*list.Element),
		queue:		list.New(),
	}
}

func (c *LRUCache) Get(key string) (int, bool) {
	elem, exists := c.items[key]
	if !exists {
		return 0, false
	}

	// 1. Двигаем элемент в начало очереди	
	c.queue.MoveToFront(elem)

	// 2. Достаем данные через приведение типов к нашей структуре *pair
	// т.к. elem.Value возвращает пустой интерфейс any (interface{})
	p := elem.Value.(*pair)
	return p.value, true
}

func (c *LRUCache) Set(key string, value int) {
	if elem, exists := c.items[key]; exists {
		p := elem.Value.(*pair)
		p.value = value
		c.queue.MoveToFront(elem) // Стал самым свежим
		return
	}

	//новый ключ
	if c.queue.Len() >= c.capacity {
		// 1. Берем самый старый элемент из хвоста списка
		oldest := c.queue.Back()
		if oldest != nil {
			oldestPair := oldest.Value.(*pair)
			delete(c.items, oldestPair.key)
			c.queue.Remove(oldest)
		}
	}

	// Создаем новую пару и пушим в начало списка
	newPair := &pair{key: key, value: value}
	newElement := c.queue.PushFront(newPair)

	// Сохраняем указатель на элемент списка в мапу
	c.items[key] = newElement
}

func main() {
	// Создаем кэш всего на 2 элемента
	cache := NewLRUCache(2)

	cache.Set("A", 1)
	cache.Set("B", 2)
	
	// Читаем "A". Теперь "A" — самый свежий, а "B" подвинулся к выходу
	_, _ = cache.Get("A") 

	fmt.Println("--- Добавляем третий элемент 'C' (Лимит превышен) ---")
	cache.Set("C", 3) // "B" должен вылететь!

	// Проверяем, кто остался
	_, existsA := cache.Get("A")
	_, existsB := cache.Get("B")
	_, existsC := cache.Get("C")

	fmt.Printf("A в кэше: %t (ожидаем true)\n", existsA)
	fmt.Printf("B в кэше: %t (ожидаем false)\n", existsB)
	fmt.Printf("C в кэше: %t (ожидаем true)\n", existsC)
}
