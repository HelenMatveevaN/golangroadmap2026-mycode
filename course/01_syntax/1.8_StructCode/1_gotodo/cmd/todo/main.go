package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gotodo/internal/storage"
	"gotodo/internal/task"	
)

/*
go run ./cmd/todo/ add Пойти выпить кофе --priority high
go run ./cmd/todo/ list
go run ./cmd/todo/ done 1

# 1. Запуск стандартного линтера Go
go vet ./...

# 2. Запуск продвинутого staticcheck
staticcheck ./...
or
golangci-lint run ./...

*/

func main() {
	// 1. Инициализируем хранилище (~/.gotodo/tasks.json)
	store, err := storage.NewStore()
	if err != nil {
		fmt.Printf("Ошибка инициализации хранилища: %v\n", err)
		os.Exit(1)
	}

	// 2. Инициализируем сервис бизнес-логики ("мозг")
	service := task.NewService(store)

	// Проверяем, передана ли команда
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		priority := addCmd.String("priority", "low", "Приоритет задачи (low, medium, high)")
		due := addCmd.String("due", "", "Дедлайн задачи (например, 2d, 12h)")

		// Хитрый трюк: отделяем флаги от текста задачи, чтобы Go их правильно распарсил
		var flagArgs []string
		var titleArgs []string

		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if strings.HasPrefix(arg, "-") {
				flagArgs = append(flagArgs, arg)
				// Если у флага есть значение (следующий аргумент), забираем и его
				if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
					flagArgs = append(flagArgs, os.Args[i+1])
					i++
				}
			} else {
				titleArgs = append(titleArgs, arg)
			}
		}

		// Парсим только чистые флаги
		_ = addCmd.Parse(flagArgs)

		if len(titleArgs) < 1 {
			fmt.Println("Ошибка: укажите название задачи. Пример: gotodo add \"Купить хлеб\" --priority high")
			os.Exit(1)
		}
		title := strings.Join(titleArgs, " ")

		// Передаем service вместо старого store
		if err := AddTask(service, title, *priority, *due); err != nil {
			fmt.Printf("Ошибка при добавлении задачи: %v\n", err)
			os.Exit(1)
		}

	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		filter := listCmd.String("filter", "all", "Фильтр задач (all, active, done)")

		_ = listCmd.Parse(os.Args[2:])

		// Передаем service и распакованную строку фильтра (*filter)
		if err := ListTasks(service, *filter); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			os.Exit(1)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: укажите ID задачи. Пример: gotodo done 1")
			os.Exit(1)
		}
		idStr := os.Args[2]
		if err := CompleteTask(service, idStr); err != nil {
			fmt.Printf("Ошибка при выполнении задачи: %v\n", err)
			os.Exit(1)
		}

	case "rm":
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: укажите ID для удаления. Пример: gotodo rm 1")
			os.Exit(1)
		}
		idStr := os.Args[2]
		if err := RemoveTask(service, idStr); err != nil {
			fmt.Printf("Ошибка при удалении задачи: %v\n", err)
			os.Exit(1)
		}

	case "clear":
		if err := ClearTasks(service); err != nil {
			fmt.Printf("Ошибка при очистке задач: %v\n", err)
			os.Exit(1)
		}

	case "help", "-h", "--help":
		printHelp()

	default:
		fmt.Printf("Неизвестная команда: %q\n", command)
		os.Exit(1)
	}
}

// printHelp выводит подсказку по использованию утилиты
func printHelp() {
	fmt.Println("CLI-менеджер задач gotodo")
	fmt.Println("\nИспользование:")
	fmt.Println("  gotodo <команда> [аргументы] [флаги]")
	fmt.Println("\nДоступные команды:")
	fmt.Println("  add \"текст\"   Добавить задачу")
	fmt.Println("                Флаги: --priority (low|medium|high), --due (например: 2d, 5h)")
	fmt.Println("  list           Показать задачи")
	fmt.Println("                Флаги: --filter (all|active|done)")
	fmt.Println("  done <id>      Отметить задачу как выполненную")
	fmt.Println("  rm <id>        Удалить задачу по ID")
	fmt.Println("  clear          Очистить базу данных")
}
