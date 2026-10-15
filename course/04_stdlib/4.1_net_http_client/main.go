package main

import (
	//"context"
	"fmt"
	"math"
	"net/http"
	"time"
)

// 1. Обертка над стандартным клиентом
type SimpleRetryClient struct {
	client	    *http.Client  // Стандартный клиент Go, который делает реальную работу
	maxRetries  int 		  // сколько раз пытаться после первой неудачи (например, 3)
	minBackoff  time.Duration // начальная пауза (например, 100ms)
}

func NewSimpleRetryClient(client *http.Client, maxRetries int, minBackoff time.Duration) *SimpleRetryClient {
	return &SimpleRetryClient{
		client:     client,
		maxRetries: maxRetries,
		minBackoff: minBackoff,
	}
}

func (c *SimpleRetryClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// Цикл идет от 0 до maxRetries (всего попыток: maxRetries + 1)
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		
		// Если это повторная попытка — спим перед ней
		if attempt > 0 {
			// Формула экспоненты: minBackoff * 2^(attempt-1)
			// При minBackoff = 100ms: попытка 1 = 100ms, попытка 2 = 200ms, попытка 3 = 400ms
			backoff := c.minBackoff * time.Duration(math.Pow(2, float64(attempt-1)))

			// Спим аккуратно. Если пока мы спим, контекст приложения отменится — мы сразу выйдем
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(backoff):
			}
		}

		// Если у запроса есть тело, сбрасываем его в начало перед повторным чтением
		if req.Body != nil && req.GetBody != nil {
			body, _ := req.GetBody()
			req.Body = body
		}

		// Делаем сам запрос
		resp, err = c.client.Do(req)

		// Проверяем: если всё хорошо (ошибки нет и статус < 500) — прерываем цикл и возвращаем ответ
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// Если мы решили идти на следующий круг ретрая, но сервер вернул нам ответ (например, 503),
		// нужно закрыть тело старого ответа, иначе сокет «залипнет»
		if resp != nil {
			_ = resp.Body.Close()
		}
	}

	// Если цикл закончился, а мы здесь — значит, все попытки исчерпаны
	if err != nil {
		return nil, fmt.Errorf("все попытки исчерпаны, последняя сетевая ошибка: %w", err)
	}
	return nil, fmt.Errorf("сервер вернул ошибку после всех ретраев, код: %d", resp.StatusCode)
}

func main() {
	// 1. Создаем обычный http.Client
	baseClient := &http.Client{Timeout: 5 * time.Second}

	// 2. Оборачиваем в наш простой ретраер
	retryClient := NewSimpleRetryClient(baseClient, 3, 100*time.Millisecond)

	// 3. Создаем запрос
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)

	// 4. Выполняем. Если сервис упадет, код сам повторит запрос до 3 раз.
	resp, err := retryClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Успешный ответ:", resp.StatusCode)
}