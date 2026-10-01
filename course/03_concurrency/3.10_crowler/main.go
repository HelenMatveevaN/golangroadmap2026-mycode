package main

import (
	"fmt"
	"context"
	"os"
	"os/signal"
	"syscall"
	"sync"
	"time"
)

// Потокобезопасный кэш адресов
type CrawCache struct {
	mu			sync.Mutex
	visited 	map[string]bool //[URL]посещен или нет
}

// Проверка: есть ли URL в карте?
func (c *CrawCache) CheckAndSet(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.visited[url]; !exists {
		c.visited[url] = true
		return false
	}
	
	return true //есть в карте (т.е., уже посещен)
}

// Мок-функций скачивания страницы (Fetcher)
func fetch(ctx context.Context, url string) ([]string, error) {
	select {
	case <-time.After(3 * time.Second):
		//продолжаем работу
		if url == "https://golang.org" {
			return []string{"https://golang.org"}, nil
		}
		return []string{}, nil

	case <-ctx.Done():
		fmt.Printf("[Fetcher] Скачивание %s прервано по Ctrl+C\n", url)
		return []string{}, ctx.Err()
	}
}

// Рекурсивная конкурентная функция обхода
func Crawl(ctx context.Context, url string, depth int, cache *CrawCache, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default: //идем дальше
	}

	//проверка глубины
	if depth <= 0 {
		return
	}

	//проверка дубликатов
	if cache.CheckAndSet(url) {
		return
	}

	sem <- struct{}{} 					//занять слот
	childUrls, err := fetch(ctx, url)   //скачивание
	<-sem       					    //освободить слот
	if err != nil {
		return
	}

	//Конкуретный запуск
	for _, childUrl := range childUrls {
		wg.Add(1)
		go Crawl(ctx, childUrl, depth-1, cache, sem, wg)
	}
}

func main() {
	var wg sync.WaitGroup

	sem := make(chan struct{}, 2) //канал-семафор, лимит 10
	
	cache := &CrawCache{visited: make(map[string]bool)}

	//родительский контекст
	parentCtx, cancelParent := context.WithCancel(context.Background())
	defer cancelParent()

	//контекст по таймауту
	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelTimeout()

	//сигнальный контекст
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	//горутина-мост (объединяет 2 контекста)
	go func() {
		select{
		case <-timeoutCtx.Done():
		case <-signalCtx.Done():
		}
		cancelParent()
	}()

	urls := []string{
		"https://go.dev",
		"https://golang.org",
		"https://habr.com",
		"https://github.com",
		"https://stackoverflow.com",
		"https://wikipedia.org",
	}	

	for _, url := range urls {
		//запуск url
		wg.Add(1)
		go Crawl(parentCtx, url, 3, cache, sem, &wg)
	}

	wg.Wait()

	for url := range cache.visited {
		fmt.Println(url)
	}
}