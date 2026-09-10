package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	
	"gowiki/pkg/config"
	"gowiki/internal/storage"
	"gowiki/internal/wiki"
)

/*
go run ./cmd/wiki/ new MyFuturisticBlog --author "iMac Developer"
go run ./cmd/wiki/ build
go run ./cmd/wiki/ serve

go vet ./...
golangci-lint run ./...
*/

func main() {
	// Загружаем конфигурацию из файла YAML
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	//инициализируем слои
	store := storage.NewDiskStore()
	wikiService = wiki.NewService(store)

	// Проверяем, передана ли команда
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "serve":
		//Запустить интерактивный веб-сервер
		http.HandleFunc("/view/", viewHandler)
		http.HandleFunc("/edit/", editHandler)
		http.HandleFunc("/save/", saveHandler)

		fmt.Printf("🌐 Режим SERVE: сервер запущен на http://localhost%s/view/TestPage", cfg.Server.Port)
		log.Fatal(http.ListenAndServe(cfg.Server.Port, nil))

	case "new":
		// Создание черновика статьи
		newCmd := flag.NewFlagSet("new", flag.ExitOnError)
		author := newCmd.String("author", "Anonymous", "Автор статьи") //указ-ль на строку
		_ = newCmd.Parse(os.Args[2:])

		args := newCmd.Args()
		if len(args) <= 2 {
			fmt.Println("Ошибка: укажите название страницы. Пример: go run ./cmd/wiki/ new MyNewPage")
			os.Exit(1)
		}
		title := args[0]

		err := wikiService.CreateNewPage(title, *author)
		if err != nil {
			log.Fatalf("Не удалось создать страницу: %v", err)
		}

		fmt.Printf("📝 Успешно создан черновик статьи: %s.md\n", title)

	case "build":
		// Генерация статического HTML сайта в папку public/
		fmt.Println("📦 Режим BUILD: Старт генерации статического сайта...")


		// Создаем папку public, если её нет
		err := os.MkdirAll("public", 0755)
		if err != nil {
			log.Fatalf("Не удалось создать папку public: %v", err)
		}

		// Получаем список всех статей
		titles, err := wikiService.GetAllTitles()
		if err != nil {
			log.Fatalf("Не удалось получить список статей: %v", err)
		}

		// Парсим шаблон отображения
		tmpl, err := template.ParseFiles("cmd/wiki/view.html")
		if err != nil {
			log.Fatalf("Ошибка шаблона view.html: %v", err)
		}

		for _, title := range titles {
			p, err := wikiService.GetPage(title)
			if err != nil {
				fmt.Printf("Пропуск статьи %s из-за ошибки: %v\n", title, err)
				continue
			}

			// Создаем физический .html файл в папке public/
			htmlFile, err := os.Create(fmt.Sprintf("public/%s.html", title))
			if err != nil {
				fmt.Printf("Не удалось создать файл для %s: %v\n", title, err)
				continue
			}

			// Рендерим данные страницы прямо в этот файл!
			_ = tmpl.Execute(htmlFile, p)
			htmlFile.Close()
			fmt.Printf("  └─ Generated: public/%s.html\n", title)
		}
		fmt.Println("🎉 Сборка успешно завершена! Проверьте папку public/")

	default:
		fmt.Printf("Неизвестная команда: %q\n", command)
		printHelp()
		os.Exit(1)	
	}
}

func printHelp() {
	fmt.Println("Использование: go run ./cmd/wiki/ <команда> [флаги]")
	fmt.Println("Команды:")
	fmt.Println("  serve   Запустить интерактивный веб-сервер")
	fmt.Println("  new     Создать новую страницу (Пример: new MyPage --author \"Ivan\")")
	fmt.Println("  build   Скомпилировать все Markdown файлы в статический HTML в папку public/")
}