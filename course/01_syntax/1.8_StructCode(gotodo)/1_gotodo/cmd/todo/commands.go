//commands.go — это официант (интерфейс)
package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"gotodo/internal/task"
)

// AddTask добавляет новую задачу
func AddTask(service *task.Service, title, priority, dueStr string) error {
	var dueDate time.Time
	var err error

	if dueStr != "" {
		// Переводим строковую дату в объект time.Time
		dueDate, err = time.Parse("2006-01-02", dueStr)
		if err != nil {
			return fmt.Errorf("неверный формат даты (нужен YYYY-MM-DD): %w", err)
		}
	} else {
			// Дедлайн по умолчанию — плюс 24 часа
			dueDate = time.Now().Add(24 * time.Hour)
	}

	// Передаем чистые данные в сервис бизнес-логики
	err = service.Create(title, priority, dueDate)
	if err != nil {
		return fmt.Errorf("не удалось добавить задачу: %w", err)
	}

	fmt.Println("Задача успешно добавлена!")
	return nil

}


// ListTasks выводит задачи с гарантированно ровными колонками на любом языке.
// Получает отфильтрованный срез из сервиса и красиво форматирует его в консоли.
func ListTasks(service *task.Service, filter string) error {
	tasks, err := service.FilteredList(filter)
	if err != nil {
		return fmt.Errorf("не удалось загрузить задачи: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Printf("Нет задач, соответствующих фильтру %q.\n", filter)
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

	for _, t := range tasks {
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
		switch strings.ToLower(t.Priority) { // Добавили strings.ToLower для надежности
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

	return nil
}

// CompleteTask отмечает задачу как выполненную
// Конвертирует строку ID в число и передает в сервис.
func CompleteTask(service *task.Service, idStr string) error {
	if idStr == "" {
		return fmt.Errorf("не указан ID задачи для удаления")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("ID задачи должен быть числом: %w", err)
	}

	err = service.Complete(id)
	if err != nil {
		return fmt.Errorf("не удалось выполнить задачу: %w", err)
	}

	fmt.Printf("Задача №%d отмечена как выполненная!\n", id)
	return nil
}

// RemoveTask удаляет задачу по ID
func RemoveTask(service *task.Service, idStr string) error {
	if idStr == "" {
		return fmt.Errorf("не указан ID задачи для удаления")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("ID задачи должен быть числом: %w", err)
	}

	err = service.Remove(id)
	if err != nil {
		return fmt.Errorf("не удалось удалить задачу: %w", err)
	}

	fmt.Printf("Задача №%d успешно удалена!\n", id)
	return nil
}

// ClearTasks полностью очищает файл
func ClearTasks(service *task.Service) error {
	err := service.Clear()
	if err != nil {
		return fmt.Errorf("не удалось очистить задачи: %w", err)
	}

	fmt.Println("Все задачи успешно удалены (файл очищен).")
	return nil
}
