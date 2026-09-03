package main

import (
	"errors"
	"fmt"
)

//Обернуть ошибку через fmt.Errorf("... %w", err).

// Инициализируем sentinel-ошибку (ошибку-маркер)
var ErrDatabaseDown = errors.New("database connection refused")

// имитация уровня доступа к данным (repository)
func fetchUserFromDB() error {
	return ErrDatabaseDown //возврат ориг.ошибки
}

// имитация бизнес-логики (service)
func getUserService() error {
	err := fetchUserFromDB()
	if err != nil {
		return fmt.Errorf("failed to get user profiles: %w", err)
	}
	return nil
}

func main() {
	err := getUserService()
	if err != nil {
		fmt.Printf("Полный лог ошибки: %v\n\n", err)

		if errors.Is(err, ErrDatabaseDown) {
			fmt.Println("🚨 Реакция системы: База данных недоступна! Включаем резервный режим.")
		} else {
			fmt.Println("Произошла какая-то другая ошибка.")
		}
	}
}