//jsonstore.go — это кладовщик (хранилище)
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gotodo/internal/task"
)

// Store управляет сохранением и загрузкой задач в JSON-файл.
type Store struct {
	filePath string
}

// NewStorage инициализирует путь к ~/.gotodo/tasks.json и создает папку при необходимости
func NewStore() (*Store, error) {
	dir, err := getStoragePath()
	if err != nil {
		return nil, fmt.Errorf("faled to get storage path: %w", err)
	}

	// Создаем папку ~/.gotodo, если её еще нет
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("faled to create directory: %w", err)
	}

	return &Store{
		filePath: filepath.Join(dir, "tasks.json"),
	}, nil
}

// getStoragePath находит домашнюю папку пользователя и возвращает путь к .gotodo.
// Имя начинается с маленькой буквы, так как функция используется только внутри этого пакета.
func getStoragePath() (string, error) {
	home, err := os.UserHomeDir() // На macOS вернет /Users/imac
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".gotodo"), nil
}

// Save сохраняет список задач в JSON-файл
func (s *Store) Save(tasks []task.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	//запись байтов в файл, 0644 - права на чтение и запись
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Load загружает список задач из JSON-файла
func (s *Store) Load() ([]task.Task, error) {
	// Если файла еще нет, возвращаем пустой список без ошибки
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []task.Task{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return tasks, nil
}
