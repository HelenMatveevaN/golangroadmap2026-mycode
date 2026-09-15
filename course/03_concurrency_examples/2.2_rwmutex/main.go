package main

import (
	"fmt"
	"sync"
	"time"
)

// Cache — кэш с блокировкой чтения/записи
type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock() // эксклюзивная блокировка для записи
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock() // разделяемая блокировка для чтения
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func main() {
	cache := NewCache()
	cache.Set("user:1", "Alice")

	var wg sync.WaitGroup

	// Запускаем 10 читателей параллельно
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			if v, ok := cache.Get("user:1"); ok {
				fmt.Printf("Читатель %d прочитал: %s\n", id, v)
			}
		}(i)
	}

	wg.Wait()
}