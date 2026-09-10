package store

import (
	"encoding/json"
	"fmt"
	"os"

	"todoapp/internal/task"
)

// ==========================================
// 1. Интерфейс TaskStore (Контракт)
// ==========================================
type TaskStore interface {
	Save(tasks []task.Task) error
	Load() ([]task.Task, error)
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

func (j *JSONStore) Save(tasks []task.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}
	if err := os.WriteFile(j.fileName, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func (j *JSONStore) Load() ([]task.Task, error) {
	if _, err := os.Stat(j.fileName); os.IsNotExist(err) {
		return []task.Task{}, nil
	}

	data, err := os.ReadFile(j.fileName)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshal tasks: %w", err)
	}
	return tasks, nil
}

// ==========================================
// 3. Реализация InMemoryStore (Хранение в памяти)
// ==========================================
type InMemoryStore struct {
	tasks []task.Task
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{tasks: make([]task.Task, 0)}
}

func (m *InMemoryStore) Save(tasks []task.Task) error {
	m.tasks = append([]task.Task(nil), tasks...) // Глубокое копирование слайса
	return nil
}

func (m *InMemoryStore) Load() ([]task.Task, error) {
	return m.tasks, nil
}

// ==========================================
// 5. Декоратор: LoggingTaskStore через эмбеддинг
// ==========================================
type LoggingTaskStore struct {
 	// LoggingTaskStore автоматически удовлетворяет контракту
	// TaskStore и получает все его методы.
	TaskStore
}

// NewLoggingTaskStore принимает любое хранилище и возвращает его обертку с логированием
func NewLoggingTaskStore(wrapped TaskStore) LoggingTaskStore {
	return LoggingTaskStore{TaskStore: wrapped}
}

// Save переопределяет оригинальный метод Save для добавления логов
func (l *LoggingTaskStore) Save(tasks []task.Task) error {
	fmt.Printf("[LOG] Попытка сохранить %d задач на диск...\n", len(tasks))
	
	// Вызываем метод исходного (обернутого) хранилища
	err := l.TaskStore.Save(tasks)

	if err != nil {
		fmt.Printf("[LOG] [ОШИБКА] Не удалось сохранить задачи: %v\n", err)
	} else {
		fmt.Println("[LOG] Данные успешно сброшены в персистентное хранилище.")
	}
	return err
}

// Load переопределяет оригинальный метод Load
func (l *LoggingTaskStore) Load() ([]task.Task, error) {
	fmt.Println("[LOG] Инициализация загрузки данных...")

	tasks, err := l.TaskStore.Load()
	if err != nil {
		fmt.Printf("[LOG] [ОШИБКА] Сбой при чтении данных: %v\n", err)
	} else {
		fmt.Printf("[LOG] Успешно загружено %d задач.\n", len(tasks))
	}

	return tasks, err
}