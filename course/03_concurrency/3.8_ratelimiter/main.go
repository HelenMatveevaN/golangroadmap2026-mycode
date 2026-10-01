package main

import (
	"fmt"
	"time"
)

func main() {
	emails := make(chan int, 10)
	for i:=1; i <= 10; i++ {
		emails <- i
	}
	close(emails)

	ticker := time.NewTicker(200 * time.Millisecond) //инт-л 200мс
	defer ticker.Stop()

	for email := range emails {
		<-ticker.C //ожидание с интервалом
		fmt.Printf("[Отправка] Письмо №%d успешно отправлено в <%v>.\n", email, time.Now().Format("15:04:05.000"))
	}
}