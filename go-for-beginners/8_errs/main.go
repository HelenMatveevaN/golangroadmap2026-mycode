package main

import (
    "errors"
    "fmt"
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

func main() {
    err := processUser(-5)
    if err != nil {
        fmt.Println(err)

        var valErr *ValidationError
        if errors.As(err, &valErr) {
            fmt.Printf("Проблемное поле: %s\n", valErr.Field)
        }
    }
}