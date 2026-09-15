package main

import (
	"context"
	"fmt"
	"time"
)

func fetchData(ctx context.Context, url string) (string, error) {

	// Имитируем медленный запрос
	select {
	case <-time.After(2 * time.Second):
		return "данные от " + url, nil
	case <-ctx.Done():
		return "", fmt.Errorf("запрос отменён: %w", ctx.Err())
	}
}

func main() {
	// Даём операции не более 500ms
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	result, err := fetchData(ctx, "https://example.com")
	if err != nil {
		fmt.Println("Ошибка:", err) // context deadline exceeded
		return
	}
	fmt.Println("Результат:", result)
}