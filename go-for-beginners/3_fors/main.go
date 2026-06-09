package main

import "fmt"

func main() {
    x := 42

    // if-else
    if x > 100 {
        fmt.Println("большое")
    } else if x > 10 {
        fmt.Println("среднее")
    } else {
        fmt.Println("маленькое")
    }

    // for — единственный цикл в Go
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }

    // while-стиль
    n := 1
    for n < 100 {
        n *= 2
    }
    fmt.Println(n)

    // range по строке
    for i, ch := range "Go!" {
        fmt.Printf("%d: %c\n", i, ch)
    }

    //Switch в Go не проваливается по умолчанию (не нужен break):
    day := "понедельник"
    switch day {
    case "суббота", "воскресенье":
        fmt.Println("выходной")
    default:
        fmt.Println("рабочий день")
    }

    //Задание: напиши программу, которая выводит числа от 1 до 20, 
    //заменяя кратные 3 на "Fizz", кратные 5 на "Buzz", кратные обоим — на "FizzBuzz".
    for i := 1; i <= 20; i++{
        if i%3==0 && i%5==0 {
            fmt.Println(i, "FizzBuzz")
        } else if i%3==0 {
            fmt.Println(i, "Fizz")
        } else if i%5==0 {
            fmt.Println(i, "Buzz")
        } else {
            fmt.Println(i, i)
        }
    }
}