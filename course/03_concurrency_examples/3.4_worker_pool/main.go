package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID  int
	Output string
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Имитируем работу
		time.Sleep(10 * time.Millisecond)
		result := Result{
			JobID:  job.ID,
			Output: fmt.Sprintf("воркер-%d обработал [%s]", id, job.Data),
		}
		results <- result
	}
}

func main() {
	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	// Запускаем воркеры
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// Отправляем задачи
	for i := 1; i <= numJobs; i++ {
		jobs <- Job{ID: i, Data: fmt.Sprintf("задача-%d", i)}
	}
	close(jobs) // сигнал воркерам: новых задач не будет

	// Закрываем results когда все воркеры завершились
	go func() {
		wg.Wait()
		close(results)
	}()

	// Собираем результаты
	for result := range results {
		fmt.Printf("Задача %d: %s\n", result.JobID, result.Output)
	}
}