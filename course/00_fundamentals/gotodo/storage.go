package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Task описывает расширенную структуру одной задачи
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	Priority  string    `json:"priority"` // low, medium, high
	DueDate   time.Time `json:"due_date"` // дедлайн
	CreatedAt time.Time `json:"created_at"`
}

// Storage управляет чтением и записью файла задач
type Storage struct {
	FilePath string
}

// находит домашнюю папку пользователя и создает там каталог .gotodo
func getStoragePath() (string, error) {
	homeDir, err := os.UserHomeDir() //// На macOS вернет /Users/imac
	if err != nil {
		return "", err
	}

	dirPath := filepath.Join(homeDir, ".gotodo")
	filePath := filepath.Join(dirPath, "tasks.json")

	// Создаем папку ~/.gotodo, если её нет
	if err = os.MkdirAll(dirPath, 0755); err != nil {
		return "", err
	}

	return filePath, nil
}

// NewStorage инициализирует путь к ~/.gotodo/tasks.json и создает папку при необходимости
func newStorage() (*Storage, error) {
	filePath, err := getStoragePath()
	if err != nil {
		return nil, err
	}
	return &Storage{FilePath: filePath}, nil
}

// Save сохраняет список задач в JSON-файл
func (s *Storage) Save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.FilePath, data, 0644) //запись байтов в файл, 0644 - права на чтение и запись
}

// Load загружает список задач из JSON-файла
func (s *Storage) Load() ([]Task, error) {
	file, err := os.ReadFile(s.FilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, err
	}

	if len(file) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(file, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
