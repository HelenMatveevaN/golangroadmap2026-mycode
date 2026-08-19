package main

import (
	"fmt"
	"os"
)

func main() {
	// Инициализируем хранилище (~/.gotodo/tasks.json)
	store, err := newStorage()
	if err != nil {
		fmt.Printf("Ошибка инициализации хранилища: %v\n", err)
		os.Exit(1)
	}

	// Проверяем, передана ли команда
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: укажите текст задачи. Пример: gotodo add \"Купить молоко\"")
			os.Exit(1)
		}
		title := os.Args[2]
		if err := AddTask(store, title); err != nil {
			fmt.Printf("Ошибка при добавлении задачи: %v\n", err)
			os.Exit(1)
		}

	case "list":
		if err := ListTasks(store); err != nil {
			fmt.Printf("Ошибка при чтении списка задач: %v\n", err)
			os.Exit(1)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: укажите ID задачи. Пример: gotodo done 1")
			os.Exit(1)
		}
		idStr := os.Args[2]
		if err := CompleteTask(store, idStr); err != nil {
			fmt.Printf("Ошибка при выполнении задачи: %v\n", err)
			os.Exit(1)
		}

	case "rm":
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: укажите ID для удаления. Пример: gotodo rm 1")
			os.Exit(1)
		}
		idStr := os.Args[2]
		if err := RemoveTask(store, idStr); err != nil {
			fmt.Printf("Ошибка при удалении задачи: %v\n", err)
			os.Exit(1)
		}

	case "clear":
		if err := ClearTasks(store); err != nil {
			fmt.Printf("Ошибка при очистке задач: %v\n", err)
			os.Exit(1)
		}

	case "help", "-h", "--help":
		printHelp()

	default:
		fmt.Printf("Неизвестная команда: %q\n", command)
		os.Exit(1)
	}

	/*priorityFlag := flag.String("priority", "medium", "Приоритет задачи (low, medium, high)")
	dueFlag := flag.String("due", "", "Срок выполнения задачи (формат: ГГГГ-ММ-ДД)")
	filterFlag := flag.String("filter", "all", "Фильтрация списка задач (all, done, pending)")

	flag.Parse()        // парсинг командной строки
	args := flag.Args() //парсим все после флагов

	//если команд нет, показываем, как пользоваться флагами
	if len(args) == 0 {
		fmt.Println("Использование: gotodo [флаги] [команда] [аргументы]")
		return
	}

	storagePath, err := getStoragePath()
	if err != nil {
		fmt.Printf("Ошибка определения пути: %v\n", err)
		os.Exit(1)
	}

	//читаем задачи из файла
	tasks, err := readTasks(storagePath)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	//первое оставшееся слово - наша команда
	command := args[0]

	switch command {
	case "add":
		//для добавления нужно 2 слова
		if len(args) < 2 {
			fmt.Println("Ошибка: Укажите текст задачи. Пример: gotodo add \"Купить хлеб\"")
			return
		}
		taskTitle := args[1]

		newID := 1
		if len(tasks) > 0 {
			newID = tasks[len(tasks)-1].ID + 1
		}

		newTask := Task{
			ID:       newID,
			Title:    taskTitle,
			Priority: *priorityFlag,
			DueDate:  *dueFlag,
			Done:     false,
		}

		tasks = append(tasks, newTask)

		err := writeTasks(storagePath, tasks)
		if err != nil {
			fmt.Printf("Ошибка сохранения: %v\n", err)
			return
		}
		fmt.Printf("Успех: Задача [%d] записана в JSON!\n", newID)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("Список задач пуст. Добавьте задачу через add")
			return
		}

		fmt.Println("=== СПИСОК ЗАДАЧ ИЗ JSON ===")
		for _, t := range tasks {
			if *filterFlag == "done" && !t.Done {
				continue
			}
			if *filterFlag == "pending" && t.Done {
				continue
			}

			status := "[ ]"
			if t.Done {
				status = " [V]"
			}

			fmt.Printf("%s %d: %s | Приоритет: %s | Срок: %s\n",
				status, t.ID, t.Title, t.Priority, t.DueDate)
		}

		fmt.Println("Показываем список задач...")
		fmt.Printf("Применен фильтр: %s\n", *filterFlag)

	case "done":
		if len(args) < 2 {
			fmt.Println("Ошибка: Укажите ID задачи. Пример: gotodo done 1")
			return
		}
		//ID задачи
		taskID := args[1]
		fmt.Printf("Отмечаем задачу с ID %s как выполненную\n", taskID)

	case "rm":
		if len(args) < 2 {
			fmt.Println("Ошибка: Укажите ID задачи для удаления. Пример: gotodo rm 1")
			return
		}
		//ID задачи
		taskID := args[1]
		fmt.Printf("Удаляем задачу с ID %s\n", taskID)

	case "clear":
		fmt.Println("Очищаем все задачи из списка...")

	default:
		fmt.Printf("Неизвестная команда: %s\n", command)
	}*/
}

// readTasks читает файл и превращает байты JSON в массив структур Go
/*func readTasks(path string) ([]Task, error) {
	// Проверяем, существует ли файл. Если нет — возвращаем пустой список
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Task{}, nil
	}

	// Читаем сырые байты из файла
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	// Десериализация: превращаем JSON-текст в слайс структур tasks
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func writeTasks(path string, task []Task) error {
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return err
	}
	//запись байтов в файл, 0644 - права на чтение и запись
	return os.WriteFile(path, data, 0644)
}
*/

/*
go run main.go --priority=high --due=2026-08-18 add "Исправить ошибку типов"
Добавляем задачу: "Исправить ошибку типов"
С флагами -> Приоритет: high, Срок: 2026-08-18
*/

// printHelp выводит подсказку по использованию утилиты
func printHelp() {
	fmt.Println("CLI-менеджер задач gotodo")
	fmt.Println("\nИспользование:")
	fmt.Println("  gotodo <команда> [аргументы]")
	fmt.Println("\nДоступные команды:")
	fmt.Println("  add \"текст\"   Добавить новую задачу")
	fmt.Println("  list           Показать все задачи")
	fmt.Println("  done <id>      Отметить задачу как выполненную")
	fmt.Println("  rm <id>        Удалить задачу по ID")
	fmt.Println("  clear          Удалить абсолютно все задачи")
	fmt.Println("  help           Показать эту справку")
}
