//task.go — это шеф-повар (бизнес-логика).
package task

import (
	"errors"
	"time"
)

// Ошибки бизнес-логики. Начинаются с префикса Err по Go-стилю.
var (
	ErrTaskNotFound = errors.New("task not found")
	ErrEmptyTitle = errors.New("task title cannot be empty")
)

// Task описывает, что такое задача в нашей системе
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Priority  string    `json:"priority"` // low, medium, high
	DueDate   time.Time `json:"due_date"` // дедлайн
	Done      bool      `json:"done"`
}

// Storage — контракт для сохранения данных. Принимаем интерфейсы по Дейву Чейни.
type Storage interface {
	Save(tasks []Task) error
	Load()	([]Task, error)
}

// Service — структура, которая объединяет логику и хранилище.
type Service struct {
	store Storage // Хранилище скрыто внутри (неэкспортируемое поле)
}

// NewService — конструктор. Возвращает конкретную структуру по правилу Google.
func NewService(s Storage) *Service {
	return &Service{store: s}
}

// Create проверяет правила бизнеса, вычисляет ID и отдает задачу на сохранение.
func (s *Service) Create(title, priority string, dueDate time.Time) error {
	// 1. Проверяем правило бизнеса
	if title == "" {
		return ErrEmptyTitle
	}

	// 2. Читаем текущие задачи через интерфейс
	tasks, err := s.store.Load()
	if err != nil {
		return err
	}

	// 3. Вычисляем ID для новой задачи
	nextID := 1
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	// 4. Создаем новую задачу
	newTask := Task{
		ID:        nextID,
		Title:     title,
		Priority:  priority,
		DueDate:   dueDate,
		Done:      false,
	}

	// 5. Добавляем в список и просим хранилище сохранить
	tasks = append(tasks, newTask)
	return s.store.Save(tasks)
}

// FilteredList просто загружает задачи. (Фильтрацию по желанию можно дописать тут).
func (s *Service) FilteredList(filter string) ([]Task, error) {
	allTasks, err := s.store.Load()
	if err != nil {
		return nil, err
	}

	if filter == "" || filter == "all" {
		return allTasks, nil
	}

	var filtered []Task

	for _, t := range allTasks {
		switch filter {
		case "done":
			if t.Done {
				filtered = append(filtered, t)
			}
		case "active":
			if !t.Done {
				filtered = append(filtered, t)
			}
		default:
			// Если передан неизвестный фильтр, возвращаем задачу (или можно возвращать ошибку)
			filtered = append(filtered, t)
		}
	}

	return filtered, nil
}

// Complete — логика отметки задачи выполненной
func (s *Service) Complete(id int) error {
	tasks, err := s.store.Load()
	if err != nil {
		return err
	}

	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}

	if !found {
		return ErrTaskNotFound
	}

	return s.store.Save(tasks)
}

// RemoveTask удаляет задачу по ID
func (s *Service) Remove(id int) error {
	tasks, err := s.store.Load()
	if err != nil {
		return err
	}

	var updatedTasks []Task
	found := false

	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue // пропускаем (удаляем) эту задачу
		}
		updatedTasks = append(updatedTasks, t)
	}

	if !found {
		return ErrTaskNotFound
	}

	return s.store.Save(updatedTasks)
}

// Clear — логика полной очистки
func (s *Service) Clear() error {
	return s.store.Save([]Task{})
}