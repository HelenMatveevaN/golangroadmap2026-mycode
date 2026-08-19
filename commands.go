package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

// AddTask добавляет новую задачу
func AddTask(s *Storage, title string) error {
	if title == "" {
		return fmt.Errorf("название задачи не может быть пустым")
	}

	tasks, err := s.Load()
	if err != nil {
		return err
	}

	//get newID
	newID := 1
	if len(tasks) > 0 {
		maxID := tasks[0].ID
		for _, t := range tasks {
			if t.ID > maxID {
				maxID = t.ID
			}
		}
		newID = maxID + 1
	}

	task := Task{
		ID:        newID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}

	tasks = append(tasks, task)
	if err := s.Save(tasks); err != nil {
		return err
	}

	fmt.Printf("Задача №%d успешно добавлена\n", newID)
	return nil
}

// ListTasks выводит все задачи в виде цветной таблицы
func ListTasks(s *Storage) error {
	tasks, err := s.Load()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("Список задач пуст.")
		return nil
	}

	// Создаем функции для раскраски текста
	blue := color.New(color.BgBlue, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	// Инициализируем tabwriter для красивых колонок в терминале
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Подсвечиваем заголовки синим
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", blue("ID"), blue("СТАТУС"), blue("ЗАДАЧА"), blue("СОЗДАНА"))

	for _, t := range tasks {
		var status string
		if t.Done {
			status = green("[x]") //зеленый для выполненных
		} else {
			status = yellow("[ ]") //желтый для активных
		}
		// Форматируем дату в удобный вид
		dateStr := t.CreatedAt.Format("02.01.2006 15:04")
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, status, t.Title, dateStr)
	}

	return w.Flush()
}

// CompleteTask отмечает задачу как выполненную
func CompleteTask(s *Storage, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("неверный формат ID, введите число")
	}

	tasks, err := s.Load()
	if err != nil {
		return err
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
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	if err := s.Save(tasks); err != nil {
		return err
	}

	fmt.Printf("Задача №%d отмечена как выполненная", id)
	return nil
}

// RemoveTask удаляет задачу по ID
func RemoveTask(s *Storage, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("неверный формат ID, введите число")
	}

	tasks, err := s.Load()
	if err != nil {
		return err
	}

	foundIdx := -1
	for i, t := range tasks {
		if t.ID == id {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	// Удаляем элемент из среза
	tasks = append(tasks[:foundIdx], tasks[foundIdx+1:]...)

	if err := s.Save(tasks); err != nil {
		return err
	}

	fmt.Printf("Задача №%d успешно удалена!\n", id)
	return nil
}

// ClearTasks полностью очищает файл
func ClearTasks(s *Storage) error {
	// Сохраняем пустой срез задач
	if err := s.Save([]Task{}); err != nil {
		return err
	}
	fmt.Println("Все задачи успешно удалены!")
	return nil
}
