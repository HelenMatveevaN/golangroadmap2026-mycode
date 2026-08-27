package main

import (
	"fmt"
	"strings"
)

func CustomLog(prefix string, args ...any) {
	var strArgs []string

	for _, arg := range args {
		strArgs = append(strArgs, fmt.Sprintf("%v", arg))
	}

	message := strings.Join(strArgs, " ")

	fmt.Printf("[%s] %s\n", prefix, message)
}

func main() {
	CustomLog("INFO", "Сервер", "успешно", "запущен")

	userID := 42
	isPremium := true
	CustomLog("DEBUG", "Пользователь:", userID, "| Премиум:", isPremium)

	extraDetails := []any{"Ошибка", 500, "Internal Server Error"}
	CustomLog("ERROR", extraDetails...)
}