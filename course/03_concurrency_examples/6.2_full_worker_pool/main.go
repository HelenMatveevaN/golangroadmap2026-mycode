package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Task struct {
	ID      int
	Payload string
}

type WorkerPool struct {
	numWorkers int
	tasks      chan Task
	processed  atomic.Int64
	wg         sync.WaitGroup
}

func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: workers,
		tasks:      make(chan Task, queueSize),
	}
}

func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
}

func (p *WorkerPool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	fmt.Printf("[воркер %d] запущен\n", id)

	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				fmt.Printf("[воркер %d] канал задач закрыт, завершаю\n", id)
				return
			}
			p.processTask(ctx, id, task)
		case <-ctx.Done():
			fmt.Printf("[воркер %d] получен сигнал отмены: %v\n", id, ctx.Err())
			return
		}
	}
}

func (p *WorkerPool) processTask(ctx context.Context, workerID int, task Task) {
	select {
	case <-ctx.Done():
		fmt.Printf("[воркер %d] пропускаем задачу %d — контекст отменён\n", workerID, task.ID)
		return
	default:
	}

	time.Sleep(50 * time.Millisecond)
	p.processed.Add(1)
	fmt.Printf("[воркер %d] задача %d выполнена: %s\n", workerID, task.ID, task.Payload)
}

func (p *WorkerPool) Submit(task Task) bool {
	select {
	case p.tasks <- task:
		return true
	default:
		return false // очередь заполнена
	}
}

func (p *WorkerPool) Shutdown() {
	close(p.tasks)
	p.wg.Wait()
	fmt.Printf("\nВсего обработано задач: %d\n", p.processed.Load())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := NewWorkerPool(5, 100)
	pool.Start(ctx)

	for i := 1; i <= 20; i++ {
		pool.Submit(Task{
			ID:      i,
			Payload: fmt.Sprintf("обработать данные #%d", i),
		})
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Для примера — останавливаем через 500ms
	go func() {
		time.Sleep(500 * time.Millisecond)
		quit <- syscall.SIGTERM
	}()

	<-quit
	fmt.Println("\nОстанавливаем пул воркеров...")
	cancel()
	pool.Shutdown()
}