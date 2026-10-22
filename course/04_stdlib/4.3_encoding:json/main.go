package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type RussianDate time.Time

//кастомный маршаллинг (в JSON), дата в формате dd.MM.yyyy
func (rd RussianDate) MarshalJSON() ([]byte, error) {
	t := time.Time(rd)
	formatted := t.Format("02.01.2006")

	jsonString := fmt.Sprintf(`"%s"`, formatted)

	return []byte(jsonString), nil
}

type Invoice struct {
	Number int 			`json:"number"`
	Date   RussianDate	`json:"date"`
}

func main() {
	inv := Invoice{
		Number: 777,
		Date:   RussianDate(time.Now()),
	}

	jsonData, _ := json.Marshal(inv)
	fmt.Println("Сериализованный JSON:")
	fmt.Println(string(jsonData))
	// Выведет: {"number":777,"date":"19.10.2026"}
}