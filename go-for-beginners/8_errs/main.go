//Модуль 8. Обработка ошибок

/*
Три правила обработки ошибок:

- Проверяй ошибку сразу после вызова функции
- Оборачивай с fmt.Errorf("%w", err), чтобы не терять контекст
- Используй errors.As / errors.Is для проверки конкретных типов
*/

//Задание: напиши парсер конфига из строки вида key=value 
//с валидацией и собственным типом ошибки.

package main

import (
    "errors"
    "fmt"
    "strings"
)

type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("поле %s: %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{Field: "age", Message: "не может быть отрицательным"}
    }
    if age > 150 {
        return &ValidationError{Field: "age", Message: "нереалистичное значение"}
    }
    return nil
}

func processUser(age int) error {
    if err := validateAge(age); err != nil {
        return fmt.Errorf("processUser: %w", err)
    }
    return nil
}

//Задание: напиши парсер конфига из строки вида key=value 
//с валидацией и собственным типом ошибки

// ConfigError для ошибок парсинга конфигурации (наш новый тип ошибки)
type ConfigError struct {
    Input  string
    Reason string
}

func (e *ConfigError) Error() string {
    return fmt.Sprintf("ошибка конфигурации %q: %s", e.Input, e.Reason)
}

//Низкоуровневый парсер
func isValidKV(s string) (string, string, error) {
    key, value, found := strings.Cut(s, "=")

    if !found {
        return "", "", &ConfigError{Input: s, Reason: "отсутствует символ '='"}
    }
    if key == "" {
        return "", "", &ConfigError{Input: s, Reason: "пустой ключ (KEY)"}
    }
    if value == "" {
        return "", "", &ConfigError{Input: s, Reason: "пустое значение (VALUE)"}
    }
    return key, value, nil
}

//Высокоуровневый парсер, оборачивает ошибку Низкоур.парсера
func parseConfigLine(line string) (string, string, error) {
    key, val, err := isValidKV(line)
    if err != nil {
        //обертка для ConfigError, добавляем контекст
        return "", "", fmt.Errorf("parseConfigLine failed: %w", err)
    }
    return key, val, err
}

func main() {
    // 1. Проверка старой логики с возрастом
    err := processUser(-5)
    if err != nil {
        fmt.Println(err)

        var valErr *ValidationError
        if errors.As(err, &valErr) {
            fmt.Printf("Проблемное поле: %s\n", valErr.Field)
        }
    }

    // 2. Тест нового парсера конфигурации
    fmt.Println("\n--- Тест парсера конфигурации ---")

    // Примеры для проверки: "KEY=VALUE", "INVALID_LINE", "=ONLY_VAL", "ONLY_KEY="
    text := "ONLY_KEY="

    key, val, err := parseConfigLine(text)
    if err != nil {
        fmt.Println("Полный лог ошибки:", err)

        //перехват спец.ошибки паринга
        var cfgErr *ConfigError
        if errors.As(err, &cfgErr) {
            fmt.Printf("Достали из обертки -> Строка: %s, Причина: %s\n", cfgErr.Input, cfgErr.Reason)
        }
    } else {
        fmt.Printf("Успех: %s = %s\n", key, val)
    }
}