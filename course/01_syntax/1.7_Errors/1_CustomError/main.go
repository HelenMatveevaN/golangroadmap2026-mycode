package main

import (
	"errors"
	"fmt"
	"time"
)

//Создать кастомную ошибку с дополнительным контекстом, обработать через errors.As.

// Базовая ошибка (причина), которая произошла на низком уровне
var ErrDeadlock = errors.New("database deadlock detected")

// 1. Создаем кастомный тип ошибки с дополнительным контекстом
type QueryError struct {
	Query 		string
	Duration	time.Duration
	Err 		error // Исходная причина падения
}

// Реализуем интерфейс error. Вешаем метод на Pointer Receiver
func (e *QueryError) Error() string {
	return fmt.Sprintf("query [%s] failed (%v): %v", e.Query, e.Duration, e.Err)
}

// Реализуем метод Unwrap, чтобы errors.As мог искать сквозь эту структуру
func (e *QueryError) Unwrap() error {
	return e.Err
}

func ExecuteQuery() error {
	// Возвращаем указатель на нашу кастомную ошибку
	return &QueryError{
		Query:	  "SELECT * FROM users WHERE id = 42 FOR UPDATE",
		Duration: 150 * time.Millisecond,
		Err:      ErrDeadlock,		
	}
}

func main() {
	err := ExecuteQuery()
	if err != nil {
		fmt.Println("Печать ошибки как строки:")
		fmt.Println(err.Error())
	}

	fmt.Println("\n--- Обработка через errors.As ---")

	var qErr *QueryError //qErr == nil.

	if errors.As(err, &qErr) {
			fmt.Println("✅ Успешно извлекли кастомную ошибку!")
			fmt.Printf("Тяжелый запрос: %s\n", qErr.Query)
			fmt.Printf("Время выполнения: %v\n", qErr.Duration)
		} else {
			fmt.Println("❌ Ошибка другого типа")
		}	

}	