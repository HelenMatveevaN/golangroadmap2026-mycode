package main

import (
	"bytes"
	"context"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAll (t *testing.T) {
	//Создаем тестовый сервер
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/slow":
			time.Sleep(60 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		case "/error500":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello")) //	5 байт
		}
	}))
	defer mockServer.Close()

	//подмена клиента на тестового
	httpClient = mockServer.Client() 

	//структура тест-кейса
	type testCase struct {
		name			string
		urls 			[]string
		timeout			time.Duration
		wantInOutput	[]string //список подстрок, искомых в финальном выводе
	}

	//массив (таблица) со всеми сценариями
	tests := []testCase{
		{
			name:			"Сценарий 1: Успешный запрос 200 ОК",
			urls:			[]string{mockServer.URL + "/success"},
			timeout:		1 * time.Second,
			wantInOutput: 	[]string{"200", "5 bytes", "<nil>"},
		},
		{
			name:			"Сценарий 2: Сервер упал с ошибкой 500",
			urls:			[]string{mockServer.URL + "/error500"},
			timeout:		1 * time.Second,
			wantInOutput: 	[]string{"500", "0 bytes", "<nil>"},
		},
		{
			name:			"Сценарий 3: Ошибка сети (несуществующий домен)",
			urls:			[]string{"http://domain.test"},
			timeout:		1 * time.Second,
			wantInOutput: 	[]string{"0", "0 bytes", "ошибка выполнения запроса"},
		},
		{
			name:			"Сценарий 4: Несколько URL одновременно (Конкурентность)",
			urls:			[]string{
								mockServer.URL + "/url1",
								mockServer.URL + "/url2",
								mockServer.URL + "/error500",
							},
			timeout:		1 * time.Second,
			wantInOutput: 	[]string{"/url1", "/url2", "/error500", "200", "500"},
		},
		{
			name:			"Сценарий 5: Прерывание по таймауту",
			urls:			[]string{mockServer.URL + "/slow"},
			timeout:		10 * time.Millisecond,
			wantInOutput: 	[]string{"context deadline exceeded"},
		},
	}

	//запуск цикла тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//пересоздаем канал нужного объема
			GlobalChanURL = make(chan ResultURL, len(tt.urls)+1)

			var wg sync.WaitGroup

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			//запуск горутин в цикле
			for _, url := range tt.urls {
				wg.Add(1)
				go checkURL(ctx, url, &wg)
			}

			wg.Wait()
			close(GlobalChanURL)

			//перехватываем вывод
			var buf bytes.Buffer
			PrintResults(&buf)
			gotOutput := buf.String()

			//проверка: все ли ожидаемые строки попали в текстовый отчет
			for _, want := range tt.wantInOutput {
				if !strings.Contains(gotOutput, want) {
					t.Errorf("Провал сценария defer%q\nНе нашли в выводе обязательную строку: %q\nВесь вывод приложения:\n%s",
							tt.name, want, gotOutput,
					)
				}
			}

		})
	}
}

func TestMainFunc_EmptyArgs(t *testing.T) {
	//имитация: пользователь запустил CLI без параметров, очищаем os.Args
	os.Args = []string{"cmd"}

	//сбросим флаги
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	//запуск ф-ии main без параметров
	main()
}