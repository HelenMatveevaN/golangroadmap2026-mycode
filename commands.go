package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

// AddTask добавляет новую задачу
func AddTask(s *Storage, title, priority, dueStr string) error {
	if title == "" {
		return fmt.Errorf("название задачи не может быть пустым")
	}

	//Валидация приоритета
	priority = strings.ToLower(priority)
	if priority != "low" && priority != "medium" && priority != "high" {
		return fmt.Errorf("неверный приоритет %q; доступные варианты: low, medium, high", priority)
	}

	// Парсинг дедлайна (поддерживает форматы вроде "24h", "2d")
	var dueDate time.Time
	if dueStr != "" {
		// Поддержка дней (Go time.ParseDuration не знает символ 'd')
		if strings.HasSuffix(dueStr, "d") {
			daysStr := strings.TrimSuffix(dueStr, "d")
			days, err := strconv.Atoi(daysStr)
			if err != nil {
				return fmt.Errorf("неверный формат дедлайна: %v", err)
			}
			dueDate = time.Now().AddDate(0, 0, days)
		} else {
			duration, err := time.ParseDuration(dueStr)
			if err != nil {
				return fmt.Errorf("неверный формат дедлайна (пример: 12h, 2d): %v", err)
			}
			dueDate = time.Now().Add(duration)
		}
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
		Priority:  priority,
		DueDate:   dueDate,
		CreatedAt: time.Now(),
	}

	tasks = append(tasks, task)
	if err := s.Save(tasks); err != nil {
		return err
	}

	fmt.Printf("Задача №%d успешно добавлена\n", newID)
	return nil
}

// ListTasks выводит задачи с учетом фильтрации
func ListTasks(s *Storage, filter string) error {
	tasks, err := s.Load()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("Список задач пуст.")
		return nil
	}

	// Инициализируем tabwriter для красивых колонок в терминале
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.StripEscape)

	// Подсвечиваем заголовки синим
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
		color.BlueString("ID"),
		color.BlueString("СТАТУС"),
		color.BlueString("ПРИОРИТЕТ"),
		color.BlueString("ЗАДАЧА"),
		color.BlueString("ДЕДЛАЙН"),
	)

	hasVisibleTasks := false
	for _, t := range tasks {
		// Фильтрация
		if filter == "active" && t.Done {
			continue
		}
		if filter == "done" && !t.Done {
			continue
		}
		hasVisibleTasks = true

		// Статус
		var status string
		if t.Done {
			status = color.GreenString("[x]") //зеленый для выполненных
		} else {
			status = color.YellowString("[ ]") //желтый для активных
		}

		// Раскраска приоритета
		var pStr string
		switch t.Priority {
		case "high":
			pStr = color.RedString("HIGH")
		case "medium":
			pStr = color.YellowString("MEDIUM")
		default:
			pStr = color.HiBlackString("LOW")
		}

		// Форматирование дедлайна
		dueStr := "-"
		if !t.DueDate.IsZero() {
			dueStr = t.DueDate.Format("02.01.2006 15:04")
			// Если задача не выполнена и дедлайн просрочен, подсветим его красным
			if !t.Done && time.Now().After(t.DueDate) {
				dueStr = color.RedString(dueStr + " [⚠️ ПРОСРОЧЕНО]")
			}
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", t.ID, status, pStr, t.Title, dueStr)
	}

	if !hasVisibleTasks {
		fmt.Printf("Нет задач, соответствующих фильтру %q.\n", filter)
		return nil
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
