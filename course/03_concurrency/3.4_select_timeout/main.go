package main

/*Практика: запрос к HTTP с таймаутом в 500 мс через select.*/

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	responseChan := make(chan string, 1)

	go func(){
		resp, err := http.Get("https://example.com")
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			return
		}

		responseChan <- string(body)
	}()

	select {
	case res := <-responseChan:
		fmt.Println("Ответ:", res)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("Ошибка: превышен таймаут в 500 мс.")
	}
}