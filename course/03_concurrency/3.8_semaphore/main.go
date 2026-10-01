package main

import (
	"fmt"
	"sync"
	"time"
)

/*
вам нужно скачать 10 файлов из интернета. 
Вы хотите запустить для каждого файла 
отдельную горутину (всего 10 горутин), 
но ваш провайдер или сервер разрешает качать 
не более 3 файлов одновременно, иначе забанит.
*/

func semafore(id int, wg *sync.WaitGroup, sem chan struct{}) {
	sem <- struct{}{} //занять место в семафоре (или ждем, если занято)
	fmt.Printf("Файл %d начал скачиваться.\n", id)
	
	time.Sleep(1 * time.Second)

	<- sem //освободили семафор
	defer wg.Done()
}

func main() {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 3)

	//скачивание 10ти файлов
	for i:= 0; i<10; i++ {
		wg.Add(1)
		i := i
		go semafore(i, &wg, sem)
	}

	wg.Wait()
	fmt.Println("Скачивание файлов завершено.")
}