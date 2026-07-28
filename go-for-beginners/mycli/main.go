package main

import (
	"context"
	"io"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

type ResultURL struct {
	URL     		string    		//какой адрес опрашивали
	StatusCode		int 	  		//HTTP-статус
	Duration 		time.Duration	//время ответа
	Size			int64			//Размер body
	Error			error
}

var GlobalChanURL chan ResultURL

var httpClient = http.DefaultClient

func init() {
	GlobalChanURL = make(chan ResultURL, 100)
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
	resp, err := httpClient.Do(req)
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

//выносим чтение из канала и печать таблицы в отдельную ф-ю
//принимает io.Writer - "куда писать" (т.е., os.Stdout для main или буфер для теста)
func PrintResults(out io.Writer) {
	w := tabwriter.NewWriter(out, 0, 8, 2, ' ', 0)
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
}

func main() {
	var wg sync.WaitGroup

	// регистрация флагов
	timeout := flag.Duration("timeout", 3*time.Second, "Ограничение времени ожидания") //возвр-ет указ-ль на значение
	flag.Parse() // парсинг командной строки
	args := flag.Args() //	регистрация параметров

	if len(args) == 0 {
		fmt.Println("Использование: go run main.go <url1> <url2> ...")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	for _, url := range args {
		wg.Add(1)
		go checkURL(ctx, url, &wg)
	}

	go func() {
		wg.Wait()
		close(GlobalChanURL)
	}()

	PrintResults(os.Stdout) //вызов новой ф-ии
}
