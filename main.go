/*Модуль 9. Горутины и каналы

Горутина — не поток ОС, а лёгкая сопрограмма. Запускай тысячи без проблем.

Канал — безопасный способ передать данные между горутинами. 
Помни: close вызывает отправитель, а не получатель.

Задание: реализуй пул воркеров — 
N горутин читают задачи из общего канала и обрабатывают их

*/

package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Воркер %d начал работу\n", id)
    time.Sleep(time.Millisecond * 100)
    fmt.Printf("Воркер %d завершил\n", id)
}

//NEW
func workerGet(idWorker int, tasks <-chan int, wg1 *sync.WaitGroup) {
    defer wg1.Done()
    for task := range tasks {
        fmt.Printf("Воркер %d обработал задачу %d\n", idWorker, task)
    }
}

func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    var wg sync.WaitGroup
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    wg.Wait()
    fmt.Println("Все воркеры завершили работу")

    nums := generate(2, 3, 4, 5)
    squares := square(nums)
    for n := range squares {
        fmt.Println(n)
    }

    //NEW
    fmt.Println("==============================================")
    /*Задание: реализуй пул воркеров — 
    N горутин читают задачи из общего канала и обрабатывают их
    */
    const numWorkers = 3

    var wg1 sync.WaitGroup
    tasks := generate(11,12,13,14,15,16,17,18,19,20)

    //запуск воркеров
    wg1.Add(numWorkers)
    for w:= 1; w <= numWorkers; w++ {
        go workerGet(w, tasks, &wg1)
    }
    
    wg1.Wait()
    fmt.Println("Все задачи обработаны.")
    
}

