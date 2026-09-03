package main

import (
	"errors"
	"fmt"
)

//Реализовать errors.Is для своего типа.

type AppError struct {
	Code 		string //Код ошибки
	Message 	string //Динамич.сообщ-е
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Is - кастомная реализация для errors.Is
// Она позволяет сравнивать нашу ошибку с эталоном только по полю Code
func (e *AppError) Is(target error) bool {
	// Приводим target к типу *AppError
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	// Ошибки считаются равными, если совпадают их коды
	return e.Code == t.Code
}

// Эталоны для проверки (Sentinel-ошибки нового типа)
var (
	ErrNotFound = &AppError{Code:  "NOT_FOUND"}
	ErrForbidden = &AppError{Code: "FORBIDDEN"}
)

func fetchDocument(id string) error {
	// Возвращаем ошибку с конкретным контекстом, но общим кодом "NOT_FOUND"
	return &AppError{
		Code:	 "NOT_FOUND",
		Message: fmt.Sprintf("document with ID %s missing in storage", id),
	}
}

func main() {
	err := fetchDocument("doc_123")
	wrappedErr := fmt.Errorf("api gateway error: %w", err)

	if errors.Is(wrappedErr, ErrNotFound) {
		fmt.Println("🎯 Сработало! Это ошибка типа ErrNotFound.")
		fmt.Printf("Оригинальный текст: %v\n", err)
	} else {
		fmt.Println("Другая ошибка")
	}
}