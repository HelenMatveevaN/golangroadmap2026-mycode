package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

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

// ListTasks выводит задачи с гарантированно ровными колонками на любом языке
func ListTasks(s *Storage, filter string) error {
	tasks, err := s.Load()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("Список задач пуст.")
		return nil
	}

	// Жестко заданные визуальные ширины колонок (в символах на экране)
	const (
		wID       = 5
		wStatus   = 12
		wPriority = 14
		wTask     = 30
	)

	// Вспомогательная функция выравнивания. Считает реальные символы (руны) на экране
	formatCol := func(coloredText string, cleanText string, totalWidth int) string {
		screenLength := utf8.RuneCountInString(cleanText) // Считаем символы, а не байты!
		padding := totalWidth - screenLength
		if padding < 0 {
			padding = 0
		}
		return coloredText + strings.Repeat(" ", padding)
	}

	// 1. Выводим заголовки шапки (теперь с точным подсчетом русских букв)
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	fmt.Printf("%s%s%s%s%s\n",
		formatCol(blue("ID"), "ID", wID),
		formatCol(blue("СТАТУС"), "СТАТУС", wStatus),
		formatCol(blue("ПРИОРИТЕТ"), "ПРИОРИТЕТ", wPriority),
		formatCol(blue("ЗАДАЧА"), "ЗАДАЧА", wTask),
		blue("ДЕДЛАЙН"),
	)

	hasVisibleTasks := false
	for _, t := range tasks {
		if filter == "active" && t.Done {
			continue
		}
		if filter == "done" && !t.Done {
			continue
		}
		hasVisibleTasks = true

		// Форматируем ID
		idStr := strconv.Itoa(t.ID)
		idCol := formatCol(idStr, idStr, wID)

		// Форматируем цветной статус
		var statusCol string
		if t.Done {
			statusCol = formatCol(color.GreenString("[x]"), "[x]", wStatus)
		} else {
			statusCol = formatCol(color.YellowString("[ ]"), "[ ]", wStatus)
		}

		// Форматируем цветной приоритет
		var priorityCol string
		switch t.Priority {
		case "high":
			priorityCol = formatCol(color.RedString("HIGH"), "HIGH", wPriority)
		case "medium":
			priorityCol = formatCol(color.YellowString("MEDIUM"), "MEDIUM", wPriority)
		default:
			priorityCol = formatCol(color.HiBlackString("LOW"), "LOW", wPriority)
		}

		// Безопасно ограничиваем длину названия задачи по рунам (символам)
		titleRunes := []rune(t.Title)
		displayTitle := t.Title
		if len(titleRunes) > wTask-3 {
			displayTitle = string(titleRunes[:wTask-6]) + "..."
		}
		titleCol := formatCol(displayTitle, displayTitle, wTask)

		// Форматируем дедлайн
		dueStr := "-"
		if !t.DueDate.IsZero() {
			dueStr = t.DueDate.Format("02.01.2006 15:04")
			if !t.Done && time.Now().After(t.DueDate) {
				dueStr = color.RedString(dueStr + " [⚠️ ПРОСРОЧЕНО]")
			}
		}

		// Выводим готовую, идеально собранную строку
		fmt.Printf("%s%s%s%s%s\n", idCol, statusCol, priorityCol, titleCol, dueStr)
	}

	if !hasVisibleTasks {
		fmt.Printf("Нет задач, соответствующих фильтру %q.\n", filter)
		return nil
	}

	return nil
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
