package main

import (
	"errors"
	"fmt"
)

//Создать кастомную ошибку с дополнительным контекстом, обработать через errors.As.
type RequestError struct {
	StatusCode int
	RequestID  string
	Err        error
}

//Реализуем интерфейс error
func (e *RequestError) Error() string {
	return fmt.Sprintf("request failed (id: %s, status: %d): %v", e.RequestID, e.StatusCode, e.Err)

}

// Unwrap позволяет errors.Is и errors.As заглядывать внутрь нашей ошибки
func (e *RequestError) Unwrap() error {
	return e.Err
}

// Имитация функции, которая может вернуть обернутую кастомную ошибку
func doNetworkRequest() error {
	baseErr := errors.New("connection timeout")

	//создаем кастомн.ошибку
	reqErr := &RequestError{
		StatusCode:	504,
		RequestID:	"req-abc-123",
		Err:		baseErr,
	}

	//обертка
	return fmt.Errorf("api layer error: %w", reqErr)
}

func main() {
	err := doNetworkRequest()
	if err != nil {
		var reqErr *RequestError

		if errors.As(err, &reqErr) {
			fmt.Println("🎉 Успешно перехватили RequestError через errors.As!")
			fmt.Printf("ID запроса: %s\n", reqErr.RequestID)
			fmt.Printf("Статус-код: %d\n", reqErr.StatusCode)
			fmt.Printf("Корневая причина: %v\n", reqErr.Unwrap())
		} else {
			fmt.Println("Это какая-то другая ошибка:", err)
		}

	}
}