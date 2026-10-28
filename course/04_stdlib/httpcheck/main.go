package main


import (
	"bufio"
	//"error"
	"fmt"
	"flag"
	"net/http"
	"io"
	"os"
	"sync"
	"time"
)

//✗ printf "https://google.com\nhttps://github.com\n" | go run main.go
//✗ echo -e "https://google.com\nhttps://github.com" | go run main.go

//go run main.go -f urls.txt
//cat urls.txt | go run main.go

type Config struct {
	FilePath string //знач-е из флага -f
}

type Result struct {
	URL 		string
	StatusCode  int
	Duration    time.Duration
	ContentLen  int64
	Err 		error
}

// hasStdinData возвращает true, если в программу передали данные через Pipe или редирект файла
func hasStdinData() bool {
	stat, err := os.Stdin.Stat() 
	if err != nil {
		return false
	}

	// Проверяем битовую маску: если флаг Named Pipe (FIFO) или регулярного файла выставлен,
	// значит, в Stdin есть данные
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// readLines - Функция чтения строк
func readLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r) //r - это файл или Stdin
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func checkURL(u string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()

	resp, err := http.Get(u)
	if err != nil {
		results<- Result{
			URL:		u,
			StatusCode: 0,
			Duration:   0,
			ContentLen: 0,
			Err: 		err,
		}
		return
	}
	defer resp.Body.Close()

	//считаем байты без выделения памяти
	bodySize, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		results<- Result{
			URL:		u,
			StatusCode: resp.StatusCode,
			Duration:   0,
			ContentLen: 0,
			Err: 		err,
		}
		return
	}

	duration := time.Since(start)

	results<- Result{
		URL:		u,
		StatusCode: resp.StatusCode,
		Duration:   duration,
		ContentLen: bodySize,
		Err: 		nil,
	}	
}

func main() {
	filePath := flag.String("f", "", "Путь к файлу со списком URL")
	flag.Parse()

	var urls []string
	var err error
	var input io.Reader

	if *filePath != "" { //передан флаг с файлом
		file, err := os.Open(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка чтения файла: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		input = file

	} else {
		if !hasStdinData() {
			fmt.Println("Ошибка: os.Stdin пуст. Программа ожидает данные через Pipe.")
			fmt.Println("Использование: echo \"https://google.com\" | go run main.go")
			os.Exit(1) // Мягко выходим с кодом 1, так как утилиту запустили неправильно
		}

		// 2. Если данные есть — спокойно читаем их
		fmt.Println("Считываем данные из os.Stdin...")

		//в os.Stdin есть данные
		if hasStdinData() {
			input = os.Stdin
		}
	}

	if input == nil {
		fmt.Println("Данных о файле нет")
		os.Exit(1)
	} 

	//данные о файле есть
	urls, err = readLines(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения из файла: %v\n", err)
		os.Exit(1)
	}

	if len(urls) == 0 {
		fmt.Println("Использование: утилита -f urls.txt или cat urls.txt | утилита")
		flag.Usage()
		os.Exit(1)
	}

	resultsChan := make(chan Result, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go checkURL(url, resultsChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Пример обработки полученных URL
	fmt.Printf("Успешно загружено URL: %d\n", len(urls))
	for res := range resultsChan {
		fmt.Println("Данные из канала", res)
	}
}