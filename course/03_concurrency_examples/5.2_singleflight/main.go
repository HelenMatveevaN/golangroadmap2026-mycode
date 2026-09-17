package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type CacheWithSingleflight struct {
	sf    singleflight.Group
	mu    sync.RWMutex
	cache map[string]string
}

func (c *CacheWithSingleflight) Get(key string) (string, error) {
	c.mu.RLock()
	if v, ok := c.cache[key]; ok {
		c.mu.RUnlock()
		return v, nil
	}
	c.mu.RUnlock()

	// singleflight: только одна горутина сходит в БД
	v, err, shared := c.sf.Do(key, func() (interface{}, error) {
		fmt.Printf("Идём в БД за ключом '%s'...\n", key)
		time.Sleep(100 * time.Millisecond)
		return "значение_" + key, nil
	})

	if err != nil {
		return "", err
	}

	fmt.Printf("Результат для '%s' был общим: %v\n", key, shared)

	c.mu.Lock()
	if c.cache == nil {
		c.cache = make(map[string]string)
	}
	c.cache[key] = v.(string)
	c.mu.Unlock()

	return v.(string), nil
}

func main() {
	cache := &CacheWithSingleflight{}
	var wg sync.WaitGroup

	// 10 горутин одновременно запрашивают один ключ
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			v, _ := cache.Get("user:1")
			fmt.Printf("Горутина %d получила: %s\n", id, v)
		}(i)
	}

	wg.Wait()
	// "Идём в БД..." напечатается только один раз
}