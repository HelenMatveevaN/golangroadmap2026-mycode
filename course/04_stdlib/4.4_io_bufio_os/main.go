package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

//генерация файла на 1Гб

func main() {
	startTime := time.Now()

	file, err := os.Create("words_1gb.txt")
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 64*1024)
	defer writer.Flush()

	// Набор слов для генерации текста
	words := []string{
		"apple", "banana", "cherry", "golang", "interface", "reader", "writer",
		"buffer", "scanner", "file", "system", "performance", "code", "compiler",
		"data", "stream", "concurrency", "goroutine", "channel", "mutex",
	}

	var targetSize int64 = 1 * 1024 * 1024 * 1024
	var currentSize int64 = 0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	wordsLen := len(words)

	fmt.Println("Генерация файла на 1 ГБ началась...")

	wordCount := 0 // Счетчик слов для контроля переноса строк

	// 3. Цикл записи данных
	for currentSize < targetSize {
		// Выбираем случайное слово
		word := words[r.Intn(wordsLen)]
		wordCount++

		// Каждые 12 слов делаем перенос строки (\n), иначе — просто пробел
		if wordCount%12 == 0 {
			word += "\n"
		} else {
			word += " "
		}
		
		// Пишем в буфер
		n, err := writer.WriteString(word)
		if err != nil {
			log.Fatalf("Ошибка при записи: %v", err)
		}
		
		// Увеличиваем счетчик байт
		currentSize += int64(n)
	}

	fmt.Printf("Успешно! Файл words_1gb.txt создан.\n")
	fmt.Printf("Размер: %.2f ГБ\n", float64(currentSize)/(1024*1024*1024))
	fmt.Printf("Время выполнения: %v\n", time.Since(startTime))
}