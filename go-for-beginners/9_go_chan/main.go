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

    fmt.Printf("========================================\n")

    //Задание: реализуй пул воркеров — 
    //N горутин читают задачи из общего канала и обрабатывают их

    //1-создаем каналы
    numJobs := 6
    jobs := make(chan int, numJobs)
    rslts := make(chan int, numJobs)

    //2 - пул воркеров
    numWorkers := 3
    var poolWg sync.WaitGroup

    for w:=1; w<=numWorkers; w++ {
        poolWg.Add(1)
        //каждая go-рутина-воркер читает из общего канала jobs
        go func(wrkID int) {
            defer poolWg.Done()
            for job := range jobs {
                fmt.Printf("Воркер %d взял задачу %d\n", wrkID, job)
                time.Sleep(time.Millisecond * 50) //Имитация работы
                rslts <- job * 2
            }
        }(w)
    }

    //3 - отправка задач в канал
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs) //закрываем, чтобы воркеры вышли из цикла

    //4 - ждем завершения всех воркеров и закрываем канал результатов
    go func(){
        poolWg.Wait()
        close(rslts)
    }()

    //вывод результатов
    for res := range rslts {
        fmt.Printf("Результат: %d\n", res)
    }
}