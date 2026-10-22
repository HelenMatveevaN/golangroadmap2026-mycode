package main

import (
	"fmt"
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now() // Засекаем время для теста

	src, err := os.Open("words_1gb.txt")
	if err != nil {
		log.Fatalf("Не удалось открыть файл %v", err)
	}
	defer src.Close()

	reader := bufio.NewReader(src) //гибкая обработка строк разной длины (в т.ч.длинных)
	
	var totalLines int

	//частота строк
	wordFrequency := make(map[string]int, 1000)

	fmt.Println("Обработка файла началась, пожалуйста, подождите...")

	for {
		//1 строка
		line, err := reader.ReadString('\n')

		if len(line) > 0 {
			totalLines++

			//разбиваем строку на слова
			words := strings.Fields(line)

			for _, word := range words {
				wordFrequency[word]++
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
	}

	fmt.Println("\n=== ЧАСТОТА СЛОВ ===")
	for word, count := range wordFrequency {
		fmt.Printf("%-15s: %d раз\n", word, count)		
	}

	fmt.Printf("Всего строк обработано: %d\n", totalLines)
	fmt.Printf("Время выполнения:       %v\n", time.Since(start))
}