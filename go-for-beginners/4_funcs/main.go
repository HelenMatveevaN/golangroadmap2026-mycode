package main

import (
    "fmt"
    "math"
)

// Функция с двумя возвращаемыми значениями
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("деление на ноль")
    }
    return a / b, nil
}

// Вариадическая функция
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// Функция как значение
func apply(x float64, fn func(float64) float64) float64 {
    return fn(x)
}

func minMax(nums []int) (int, int) {
    min := nums[0]
    max := nums[0]
    for i:= 1; i<len(nums); i++ {
        if nums[i] < min {
            min = nums[i]
        }
        if nums[i] > max {
            max = nums[i]
        }
    }
    return min, max
}

func main() {
    result, err := divide(10, 3)
    if err != nil {
        fmt.Println("Ошибка:", err)
    } else {
        fmt.Printf("%.2f\n", result)
    }

    fmt.Println(sum(1, 2, 3, 4, 5))

    // Анонимная функция
    double := func(x float64) float64 { return x * 2 }
    fmt.Println(apply(5, double))
    fmt.Println(apply(9, math.Sqrt))


    //Паттерн (result, error) — главный способ обработки ошибок в Go. 
    //Не исключения, а явные возвращаемые значения.

    //Задание: напиши функцию minMax(nums []int) (int, int), 
    //которая возвращает минимум и максимум из среза.
    nums := []int{2,7,12,-100,9,61,3,8,1,5}
    fmt.Println(minMax(nums))
}