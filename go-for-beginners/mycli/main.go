package main

import (
	"context"
	"io"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	"text/tabwriter"
)

type ResultURL struct {
	URL     		string    		//какой адрес опрашивали
	StatusCode		int 	  		//HTTP-статус
	Duration 		time.Duration	//время ответа
	Size			int64			//Размер body
	Error			error
}

var GlobalChanURL chan ResultURL

func init() {
	GlobalChanURL = make(chan ResultURL, 3)
}

func checkURL(ctx context.Context, url string, wg *sync.WaitGroup) {
	start := time.Now()

	defer wg.Done()

	//time.Sleep(6*time.Second)

	//создаем объект запроса с привязкой к ctx, nil - тело запроса
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) //req=*http.Request
	if err != nil {
		GlobalChanURL <- ResultURL{
			URL:			url,
			StatusCode:		0,
			Duration:		time.Since(start),
			Size:			0,
			Error:			fmt.Errorf("ошибка создания запроса: %w", err),
		}
		return
	}

	//отправляем запрос через станд.клиент
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		GlobalChanURL <- ResultURL{
			URL:			url,
			StatusCode:		0,
			Duration:		time.Since(start),
			Size:			0,
			Error:			fmt.Errorf("ошибка выполнения запроса: %w", err),
		}
		return
	}
	defer resp.Body.Close()

	//читаем тело ответа
	size, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		GlobalChanURL <- ResultURL{
			URL:			url,
			StatusCode:		resp.StatusCode,
			Duration:		time.Since(start),
			Size:			size,
			Error:			err,
		}
		return
	}

	GlobalChanURL <- ResultURL{
		URL:			url,
		StatusCode:		resp.StatusCode,
		Duration:		time.Since(start),
		Size:			size,
		Error:			nil,
	}
}

func main() {

	var wg sync.WaitGroup

	// регистрация флагов
	timeout := flag.Duration("timeout", 3*time.Second, "Ограничение времени ожидания") //возвр-ет указ-ль на значение

	// парсинг командной строки
	flag.Parse()

	//	регистрация параметров
	args := flag.Args()

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	/*// вывод флагов
	fmt.Printf("--- ФЛАГИ ---\n")
	fmt.Printf("Флаг -timeout: %v\n\n", *timeout)
	*/

	// вывод параметров
	//fmt.Printf("--- ARGS ---\n")
	
	/*fmt.Printf("Всего параметров: %d\n", len(args))
	fmt.Printf("Полный список параметров: %v\n", args)

	if len(args) > 0 {
		fmt.Printf("Первый параметр (индекс 0): %s\n", flag.Arg(0))
	}
	fmt.Printf("\n")

	if len(args) < 2 {
		fmt.Println("Использование: go run main.go <url1> <url2> ...")
		return
	}
	*/

	for _, url := range args {
		wg.Add(1)
		go checkURL(ctx, url, &wg)
	}

	go func() {
		wg.Wait()
		close(GlobalChanURL)
	}()

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "URL\tSTATUS\tDURATION\tSIZE\tERROR")

	//	чтение результатов
	for res := range GlobalChanURL {
		errStr := "<nil>"
		if res.Error != nil {
			errStr = res.Error.Error()
		}
		ms := res.Duration.Milliseconds()
		fmt.Fprintf(w, "%s\t%d\t%dms\t%d bytes\t%s\n", res.URL, res.StatusCode, ms, res.Size, errStr)
	}

	w.Flush() // сброс буфера, чтобы строки вывелись на экран

	//fmt.Println("\nКонец функции main")
}
