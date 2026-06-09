package main

import (
        "fmt"
        "regexp"
        "strings"
    )

func main() {
    // Массив — фиксированный размер
    arr := [3]int{1, 2, 3}
    fmt.Println(arr)

    // Срез — динамический, используется повсеместно
    nums := []int{10, 20, 30}
    nums = append(nums, 40, 50)
    fmt.Println(nums)
    fmt.Println(nums[1:3]) // срез [20, 30]

    // make для создания среза с нужной ёмкостью
    buf := make([]byte, 0, 64)
    buf = append(buf, "hello"...)
    fmt.Println(string(buf))

    // Карта (map)
    scores := map[string]int{
        "Алексей": 95,
        "Мария":   88,
    }
    scores["Иван"] = 72

    for name, score := range scores {
        fmt.Printf("%s: %d\n", name, score)
    }

    // Проверка наличия ключа
    val, ok := scores["Пётр"]
    if !ok {
        fmt.Println("Пётр не найден, val =", val)
    }

    //Задание: получи список слов, подсчитай частоту каждого слова и выведи результат.
    text := 
    //"Go — отличный язык. Изучать Go легко, ведь Go создан для людей!"
    "Мама мыла раму, Мама ела кораллы"

    // 1. Приводим текст к нижнему регистру
    lowerText := strings.ToLower(text)

    // 2. Удаляем пунктуацию, оставляя только буквы, цифры и пробелы
    // Регулярное выражение [^\w\s] находит все, что НЕ является буквой, цифрой или пробелом
    // Для корректной работы с кириллицей используем \p{L} (любая буква) и \p{N} (любая цифра)
    reg, _ := regexp.Compile(`[^\p{L}\p{N}\s]+`) 
    //^-отрицание, \s-любые пробельные символы, + - один или несколько раз
    cleanText := reg.ReplaceAllString(lowerText,"")

    // 3. Разбиваем очищенную строку на срез слов по пробелам
    words := strings.Fields(cleanText)

    // 4. Создаем map для подсчета частоты слов
    wordFreq := make(map[string]int)

    // 5. Заполняем карту: увеличиваем счетчик для каждого встреченного слова
    for _, word := range words {
        wordFreq[word]++
    }

    // 6. Выводим результат на экран
    fmt.Println("======================")
    fmt.Println("Частота слов в тексте:")
    for word, cnt := range wordFreq {
        fmt.Printf("%s: %d\n", word, cnt)
    }
}