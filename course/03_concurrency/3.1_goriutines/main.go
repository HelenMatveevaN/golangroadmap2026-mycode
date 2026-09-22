package main

/*
запусти 100 одновременных горутин, 
посмотри память через runtime.MemStats
*/

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	var m1, m2 runtime.MemStats

	// 1. Форсируем очистку мусора и замеряем начальную память
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Канал для удержания горутин в памяти
	blockChan := make(chan struct{})
	numGoroutines := 100_000

	for i := 0; i < numGoroutines; i++ {
		go func() {
			<-blockChan
		}()
	}

	// Даем планировщику микропаузу, чтобы он успел их аллоцировать
	time.Sleep(10 * time.Millisecond)

	// 3. Замеряем память после создания горутин
	runtime.ReadMemStats(&m2)

	// Разблокируем горутины, чтобы программа корректно завершилась
	close(blockChan)

	// Выводим результаты
	fmt.Printf("Активных горутин: %d\n", runtime.NumGoroutine())
	fmt.Printf("Память ДО (Sys):    %d КБ\n", m1.Sys/1024)
	fmt.Printf("Память ПОСЛЕ (Sys): %d КБ\n", m2.Sys/1024)
	fmt.Printf("Разница (на 100 000 шт): %d КБ\n", (m2.Sys-m1.Sys)/1024)
	fmt.Printf("Примерно на 1 горутину: %.2f КБ\n", float64(m2.Sys-m1.Sys)/1024/float64(numGoroutines))

}