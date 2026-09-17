package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

func checkURL(ctx context.Context, url string) error {
	client := &http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("создание запроса для %s: %w", url, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("запрос к %s: %w", url, err)
	}
	defer resp.Body.Close()
	fmt.Printf("%s — статус: %d\n", url, resp.StatusCode)
	return nil
}

func main() {
	urls := []string{
		"https://go.dev",
		"https://pkg.go.dev",
		"https://golang.org",
	}

	// errgroup автоматически отменяет контекст при первой ошибке
	g, ctx := errgroup.WithContext(context.Background())

	for _, url := range urls {
		url := url // захватываем переменную (Go < 1.22)
		g.Go(func() error {
			return checkURL(ctx, url)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Все URL доступны!")
}