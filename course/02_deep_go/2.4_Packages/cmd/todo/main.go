package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"todoapp/internal/store"
	"todoapp/internal/task"
)

/*
go run main.go add "Купить хлеб"
go run main.go list
go run main.go done 1
go run main.go remove 1

go build -buildvcs=false -o clitodo
GOFLAGS="-buildvcs=false" golangci-lint run
./clitodo add "Сделать домашку"
./clitodo list

*/


func printHelp() {
	fmt.Println("Использование CLI Todo:")
	fmt.Println("  todo add \"Текст задачи\"  - Добавить новую задачу")
	fmt.Println("  todo list                 - Показать все задачи")
	fmt.Println("  todo done <id>            - Отметить задачу как выполненную")
	fmt.Println("  todo remove <id>          - Удалить задачу")
}


func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	baseStore := store.NewJSONStore("todo.json")
	var appStore store.TaskStore = &store.LoggingTaskStore{
		TaskStore: baseStore,
	}

	command := os.Args[1]
	tasks, err := appStore.Load()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		return
	}

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: Укажите текст задачи.")
			return
		}
		text := os.Args[2]
		newID := 1
		if len(tasks) > 0 {
			newID = tasks[len(tasks)-1].ID + 1
		}

		tasks = append(tasks, task.Task{ID: newID, Text: text, Done: false})
		_ = appStore.Save(tasks)
		fmt.Printf("✅ Задача №%d добавлена!\n", newID)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("🎉 Список задач пуст!")
			return
		}
		fmt.Println("--- Ваши Задачи ---")
		for _, task := range tasks {
			status := "❌"
			if task.Done {
				status = "✅"
			}
			fmt.Printf("%d. [%s] %s\n", task.ID, status, task.Text)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: Укажите ID задачи.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			// Оборачиваем ошибку парсинга в нашу кастомную
			err = fmt.Errorf("%w: %v", task.ErrInvalidID, err)
			handleError(err)
			return
		}

		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Done = true
				found = true
				break
			}
		}

		if !found {
			handleError(fmt.Errorf("%w: ID %d", task.ErrTaskNotFound, id))
			return
		}

		_ = appStore.Save(tasks)
		fmt.Printf("✅ Задача №%d выполнена!\n", id)

	case "remove": // Название команды успешно обновлено
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: Укажите ID задачи. Пример: todo remove 1")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			handleError(fmt.Errorf("%w: %v", task.ErrInvalidID, err))
			return
		}

		newTasks := make([]task.Task, 0, len(tasks))
		found := false
		for _, task := range tasks {
			if task.ID == id {
				found = true
				continue
			}
			newTasks = append(newTasks, task)
		}

		if !found {
			handleError(fmt.Errorf("%w: ID %d", task.ErrTaskNotFound, id))
			return
		}

		_ = appStore.Save(newTasks)
		fmt.Printf("🗑️ Задача №%d удалена!\n", id)

	default:
		printHelp()
	}
}

// Демонстрация работы с errors.Is
func handleError(err error) {
	if errors.Is(err, task.ErrInvalidID) {
		fmt.Println("🚨 Ошибка валидации: введен некорректный формат ID. Используйте целое число.")
		return
	}
	if errors.Is(err, task.ErrTaskNotFound) {
		fmt.Printf("🚨 Ошибка бизнес-логики: %v. Проверьте список через 'todo list'.\n", err)
		return
	}
	fmt.Printf("Неизвестная ошибка: %v\n", err)
}