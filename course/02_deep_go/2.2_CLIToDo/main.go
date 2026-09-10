package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
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

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// Глобальные sentinel-ошибки для обработки через errors.Is
var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidID    = errors.New("invalid task ID")
)

// ==========================================
// 1. Интерфейс TaskStore (Контракт)
// ==========================================
type TaskStore interface {
	Save(tasks []Task) error
	Load() ([]Task, error)
}

// ==========================================
// 2. Реализация JSONStore (Работа с диском)
// ==========================================
type JSONStore struct {
	fileName string
}

func NewJSONStore(fileName string) *JSONStore {
	return &JSONStore{fileName: fileName}
}

func (j *JSONStore) Save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}
	if err := os.WriteFile(j.fileName, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func (j *JSONStore) Load() ([]Task, error) {
	if _, err := os.Stat(j.fileName); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(j.fileName)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshal tasks: %w", err)
	}
	return tasks, nil
}

// ==========================================
// 3. Реализация InMemoryStore (Хранение в памяти)
// ==========================================
type InMemoryStore struct {
	tasks []Task
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{tasks: make([]Task, 0)}
}

func (m *InMemoryStore) Save(tasks []Task) error {
	m.tasks = append([]Task(nil), tasks...) // Глубокое копирование слайса
	return nil
}

func (m *InMemoryStore) Load() ([]Task, error) {
	return m.tasks, nil
}

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

	// Инициализируем хранилище через интерфейс TaskStore.
	// Если захотим переключиться на память, заменим на `NewInMemoryStore()`
	var store TaskStore = NewJSONStore("todo.json")

	command := os.Args[1]
	tasks, err := store.Load()
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

		tasks = append(tasks, Task{ID: newID, Text: text, Done: false})
		if err := store.Save(tasks); err != nil {
			fmt.Printf("Ошибка сохранения: %v\n", err)
			return
		}
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
			err = fmt.Errorf("%w: %v", ErrInvalidID, err)
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
			handleError(fmt.Errorf("%w: ID %d", ErrTaskNotFound, id))
			return
		}

		if err := store.Save(tasks); err != nil {
			fmt.Printf("Ошибка сохранения: %v\n", err)
			return
		}
		fmt.Printf("✅ Задача №%d выполнена!\n", id)

	case "remove": // Название команды успешно обновлено
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: Укажите ID задачи. Пример: todo remove 1")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			handleError(fmt.Errorf("%w: %v", ErrInvalidID, err))
			return
		}

		newTasks := make([]Task, 0, len(tasks))
		found := false
		for _, task := range tasks {
			if task.ID == id {
				found = true
				continue
			}
			newTasks = append(newTasks, task)
		}

		if !found {
			handleError(fmt.Errorf("%w: ID %d", ErrTaskNotFound, id))
			return
		}

		if err := store.Save(newTasks); err != nil {
			fmt.Printf("Ошибка сохранения: %v\n", err)
			return
		}
		fmt.Printf("🗑️ Задача №%d удалена!\n", id)

	default:
		printHelp()
	}
}

// Демонстрация работы с errors.Is
func handleError(err error) {
	if errors.Is(err, ErrInvalidID) {
		fmt.Println("🚨 Ошибка валидации: введен некорректный формат ID. Используйте целое число.")
		return
	}
	if errors.Is(err, ErrTaskNotFound) {
		fmt.Printf("🚨 Ошибка бизнес-логики: %v. Проверьте список через 'todo list'.\n", err)
		return
	}
	fmt.Printf("Неизвестная ошибка: %v\n", err)
}