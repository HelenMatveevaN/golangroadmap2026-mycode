package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Воркер %d остановлен: %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("Воркер %d работает...\n", id)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	time.Sleep(600 * time.Millisecond)
	fmt.Println("Отменяем все горутины...")
	cancel() // сигнализируем всем горутинам об остановке

	time.Sleep(100 * time.Millisecond) // даём время завершиться
}