package main

/*напиши worker pool на N воркеров, обрабатывающих задачи из канала, 
с graceful shutdown.*/

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Println("Воркер ", id, " взял задачу ", job)
		time.Sleep(2 * time.Second)
	}
	wg.Done()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fmt.Println("[Main] Программа запущена. Нажмите Ctrl+C для проверки прерывания...")

	var wg sync.WaitGroup
	numWorkers := 3
	numJobs := 20

	jobs := make(chan int)	

	i := 1
	for i <= numWorkers {
		wg.Add(1)
		go worker(i, jobs, &wg) //запуск воркеров
		i++
	}

	jobID := 1
	Loop:
	for jobID <= numJobs {
		select {
		case jobs <- jobID:
			jobID++
		case <-ctx.Done():
			fmt.Println("Получен сигнал отмены контекста, прекращаем генерацию задач")
			break Loop
		}
	}

	close(jobs)
	wg.Wait()

	fmt.Println("Все воркеры завершили свою работу. Программа остановлена.")
}