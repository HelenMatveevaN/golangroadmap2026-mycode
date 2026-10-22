package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Name		string
	Interval	time.Duration
	LastRun		time.Time
	JobFunc		func()
}

type Scheduler struct {
	mu 		sync.Mutex
	tasks 	[]*Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make([]*Task, 0),
	}
}

func (s *Scheduler) AddTask(t *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, t)
}

func main() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	myScheduler := NewScheduler()

	myScheduler.AddTask(&Task{
		Name:		"Пинг сервера",
		Interval:	3 * time.Second,
		JobFunc:    func(){
			fmt.Println("   -> [Job 1] Сервер доступен.")
		},
	})

	myScheduler.AddTask(&Task{
		Name:		"Очистка кэша",
		Interval:	5 * time.Second,
		JobFunc:    func(){
			fmt.Println("   -> [Job 2] Кэш успешно очищен.")
		},
	})

	myScheduler.AddTask(&Task{
	    Name:     "Генерация отчета",
	    Interval: 15 * time.Second,
	    JobFunc: func() {
	        fmt.Println("   -> [Job 5] 📊 Началась сборка тяжелого отчета за месяц...")
	        time.Sleep(2 * time.Second) // Имитируем долгую работу
	        fmt.Println("   -> [Job 5] ✅ Отчет готов и отправлен на почту.")
	    },
	})	

	myScheduler.AddTask(&Task{
	    Name:     "Сбор метрик",
	    Interval: 1 * time.Second,
	    JobFunc: func() {
	        fmt.Println("   -> [Job 4] CPU: 12%, RAM: 45%")
	    },
	})

	myScheduler.AddTask(&Task{
		Name:		"Проверка диска",
		Interval:	10 * time.Second,
		JobFunc:    func(){
			fmt.Println("   -> [Job 3] Свободно: 42 ГБ. Места достаточно.")
		},
	})	

	for now := range ticker.C {
		fmt.Println("Current time: ", now)

		for _, task := range myScheduler.tasks {
			if now.Sub(task.LastRun) >= task.Interval {
				task.LastRun = now

				go func(f func()) {
					f()
				}(task.JobFunc)
			}
		}
	}
}