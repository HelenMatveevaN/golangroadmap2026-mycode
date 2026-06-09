/*
Основные типы: int, float64, string, bool, byte, rune

Нулевые значения — переменная без присвоения не мусор, а ноль: 
0 для чисел, "" для строк, false для булевых.
*/

package main

import "fmt"

func main() {
    // Три способа объявить переменную
    var name string = "Алексей"
    var age int = 25
    score := 99.5 // короткое объявление, тип выводится автоматически

    fmt.Println(name, age, score)

    // Константа
    const pi = 3.14159
    fmt.Println(pi)



    //Задание: объяви переменные для описания фильма (название, год, рейтинг) и выведи их в одну строку.
    var nameF string = "Дьявол носит Прада"
    var yearF int = 2026
    scoreF := 7.0

    fmt.Println(nameF, yearF, scoreF)

}