package main

import (
	"context"
	"fmt"
)

// Используем typedKey, чтобы избежать коллизий ключей
type contextKey string

const requestIDKey contextKey = "requestID"

func middleware(ctx context.Context) context.Context {
	// Добавляем request ID в контекст (только request-scoped данные!)
	return context.WithValue(ctx, requestIDKey, "req-abc-123")
}

func handler(ctx context.Context) {
	// Извлекаем только нужный тип
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Println("Обрабатываем запрос:", reqID)
	}
}

func main() {
	ctx := context.Background()
	ctx = middleware(ctx)
	handler(ctx)
}