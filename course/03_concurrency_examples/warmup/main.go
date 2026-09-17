package main

import (
	"context"
	"fmt"
	"time"
)

// Имитируем тяжелый запрос к Oracle
func queryOracle(ctx context.Context, resultChan chan<- string) {
	// Допустим, запрос реально занимает 500 мс
	workDuration := 500 * time.Millisecond
	timer := time.NewTimer(workDuration)
	defer timer.Stop() // Чистим память в куче по правилам Senior!

	select {
	case <-timer.C:
		resultChan <- "Успешный ответ от Oracle: 42"
	case <-ctx.Done():
		// Сюда мы попадем, если контекст отменится раньше, чем сработает таймер
		fmt.Println("[Worker] Запрос к БД отменен рантаймом, очищаем ресурсы...")
		return
	}
}

func main() {
	// Создаем контекст с таймаутом на 300 мс
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	resultChan := make(chan string, 1) // Буфер в 1 элемент, чтобы горутина не заблокировалась при записи

	go queryOracle(ctx, resultChan)

	// Главный select программы
	select {
	case res := <-resultChan:
		fmt.Println("Результат:", res)
	case <-ctx.Done():
		fmt.Println("Главный поток: Увы, превышен таймаут операции:", ctx.Err())
	}
}