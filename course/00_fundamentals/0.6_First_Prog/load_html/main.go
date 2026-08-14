package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) < 2 {
		color.Yellow("Использование: go run main.go <URL>")
		return
	}

	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		color.Red("Ошибка при отправке запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		color.Red("Сервер вернул ошибку: %d %s\n", resp.StatusCode, resp.Status)
		return
	}

	htmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("Ошибка при чтении HTML: %v\n", err)
		return
	}

	fmt.Printf("%s\n", htmlData)

	color.Green("\n▲ Страница успешно загружена! Длина HTML: %d байт.", len(htmlData))
}

/*
➜  load_html git:(main) ✗ go run main.go https://example.com
<!doctype html><html lang="en"><head><title>Example Domain</title><link rel="icon" href="data:,"><meta name="viewport" content="width=device-width, initial-scale=1"><style>body{background:#eee;width:60vw;margin:15vh auto;font-family:system-ui,sans-serif}h1{font-size:1.5em}div{opacity:0.8}a:link,a:visited{color:#348}</style></head><body><div><h1>Example Domain</h1><p>This domain is for use in documentation examples without needing permission. Avoid use in operations.</p><p><a href="https://iana.org/domains/example">Learn more</a></p></div></body></html>


elena@go-server:~/00_fundamentals/0.6_First_Prog/load_html$ go run main.go https://example1.com
Сервер вернул ошибку: 406 406 Not Acceptable
elena@go-server:~/00_fundamentals/0.6_First_Prog/load_html$ go run main.go https://example.com
<!doctype html><html lang="en"><head><title>Example Domain</title><link rel="icon" href="data:,"><meta name="viewport" content="width=device-width, initial-scale=1"><style>body{background:#eee;width:60vw;margin:15vh auto;font-family:system-ui,sans-serif}h1{font-size:1.5em}div{opacity:0.8}a:link,a:visited{color:#348}</style></head><body><div><h1>Example Domain</h1><p>This domain is for use in documentation examples without needing permission. Avoid use in operations.</p><p><a href="https://iana.org/domains/example">Learn more</a></p></div></body></html>


▲ Страница успешно загружена! Длина HTML: 559 байт.


*/