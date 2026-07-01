/*
Задание: 
напиши функцию, которая параллельно делает запросы к нескольким URL 
и возвращает результат первого успешного ответа.
*/

package main

import (
    "context"
    "fmt"
    "time"
    "net/http"
)

func slowOperation(ctx context.Context) (string, error) {
    result := make(chan string, 1)
    go func() {
        time.Sleep(200 * time.Millisecond)
        result <- "готово"
    }()

    select {
    case res := <-result:
        return res, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

//NEW
func fetchURL(ctx context.Context, url string) (string, error) {

	// Имитируем, что разные сайты отвечают с разной скоростью
	if url == "https://google.com" {
		time.Sleep(500 * time.Millisecond)
	}

	if url == "https://yandex.ru" {
		time.Sleep(200 * time.Millisecond) // Яндекс чуть медленнее
	}

	if url == "https://github.com/" {
		time.Sleep(500 * time.Millisecond) // Яндекс чуть медленнее
	}

	//resp - указатель на *http.Response
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err //если лежит сайт или нет и-нета
	}

	//выполн-ем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err //если лежит сайт или нет и-нета
	}
	defer resp.Body.Close() //от network leaks (утечек сетев.соед-ий)

	if resp.StatusCode == http.StatusOK {
		return fmt.Sprintf("Ответ от %s: Успех (%d)", url, resp.StatusCode), nil
	}
	return "", fmt.Errorf("ошибка: статус %d", resp.StatusCode)
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    res, err := slowOperation(ctx)
    if err != nil {
        fmt.Println("Таймаут:", err)
    } else {
        fmt.Println("Результат:", res)
    }

    ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel2()

    res, err = slowOperation(ctx2)
    if err != nil {
        fmt.Println("Ошибка:", err)
    } else {
        fmt.Println("Результат:", res)
    }

    /*
	Задание: 
	напиши функцию, которая параллельно делает запросы к нескольким URL 
	и возвращает результат первого успешного ответа.
	*/
	ctx3, cancel3 := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel3()

	numUrls := 3
	urls := make(chan string, numUrls)

	urlList := []string{
		"https://yandex.ru", 
		"https://google.com",
		"https://yaneznayu.ru",
		"https://github.com/",
	}

	for _, s := range urlList{
		go func(url string){
			res, err := fetchURL(ctx3, url)
			if err == nil {
				select {
				case urls <- res:
				default:
				}
			}
		}(s)
	}

	select {
	case firstUrl := <- urls:
		fmt.Println("Первый успешный ответ:", firstUrl)
	case <-ctx3.Done():
		fmt.Println("Таймаут превышен, все запросы упали с ошибкой")
	}
}