package main

import (
    "fmt"
)

func testClosure() (func() int, func() int) {
    var x int = 10
    
    //объявляем анонимную функцию (создаем "пульт")
    inc := func() int {
        x++
        return x
    }
    
    dec := func() int {
        x--
        return x
    }
    
    //Функция testClosure отдает эти два пульта в main и полностью завершает свою работу
    return inc, dec
}

func main() {
    increment, decrement := testClosure()
    
    fmt.Println(increment()) // Первый вызов
    fmt.Println(decrement()) // Второй вызов
}
