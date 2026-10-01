package main

import (
	"sync"
)

/*паттерн Горутина-Владелец*/

func main() {
	var wg sync.WaitGroup
	var wg1 sync.WaitGroup
	//var mu sync.Mutex

	c := make(chan struct{}) //потокобезопасный канал (о сигнале инкремента)
	m := make(map[string]int) //сч-к просмотра страниц сайта

	//1 фоновая горутина-менеджер
	wg1.Add(1)
	go func(){
		for range c {
			m["main_page"]++
		}
		defer wg1.Done()
	}()

	//логика горутин-клиентов
	for i :=0; i<10; i++ {
		wg.Add(1)

		go func(n int){
			c <- struct{}{}
			defer wg.Done()
		}(i)
	}

	wg.Wait() //ждем завершения всех горутин
	close(c)

	wg1.Wait() //ждем завершения горутины-менеджера
	
}