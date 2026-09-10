package main

import (
	"fmt"
	"reflect"
)

/*
 Практика: функция PrintStruct(any), которая выводит все поля и их значения 
 через рефлексию.
*/

// PrintStruct анализирует переданный объект и выводит все его поля и значения
func PrintStruct(s any) {
	// Шаг 1: Получаем reflect.Value и reflect.Type
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	// Шаг 2: Обрабатываем случай, если нам передали указатель на структуру (*Struct)
	// Мы автоматически разыменовываем его через .Elem(), как требует 3-й закон рефлексии
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// Шаг 3: Валидация. Если после разыменования это не структура — выводим ошибку
	if v.Kind() != reflect.Struct {
		fmt.Printf("Ошибка: Ожидалась структура, но получен тип %s (Kind: %s)\n", t, v.Kind())
		return
	}

	fmt.Printf("=== Анализ структуры: %s ===\n", t.Name())

	// Шаг 4: Бежим циклом по всем полям структуры
	for i := 0; i < v.NumField(); i++ { 
		fieldValue := v.Field(i) // reflect.Value поля (для получения значения)
		fieldType := t.Field(i)  // reflect.StructField поля (для получения метаданных)

		// Важный нюанс: если поле неэкспортируемое (с маленькой буквы),
		// мы не можем вызвать метод .Interface() напрямую — это вызовет панику.
		// Проверяем доступность поля для чтения:
		var displayValue any
		if fieldValue.CanInterface() {
			displayValue = fieldValue.Interface()
		} else {
			displayValue = "<неэкспортируемое поле (private)>"
		}

		// Выводим имя поля, его точный тип и значение
		fmt.Printf("Поле: %-10s | Тип: %-12s | Значение: %v\n", 
			fieldType.Name, 
			fieldType.Type, 
			displayValue,
		)
	}
	fmt.Println()
}	

// Тестовые структуры
type User struct {
	ID        int
	Username  string
	IsActive  bool
	secretKey string // Неэкспортируемое поле для проверки безопасности рефлексии
}

type Task struct {
	Title string
	Hours float64
}

func main() {
	// 1. Тестируем со структурой по значению
	u := User{ID: 42, Username: "alex_dev", IsActive: true, secretKey: "123-qwerty"}
	PrintStruct(u)

	// 2. Тестируем со структурой по указателю (должно отработать так же хорошо благодаря .Elem())
	t := &Task{Title: "Изучить рефлексию", Hours: 2.5}
	PrintStruct(t)

	// 3. Тестируем передачу некорректного типа (для проверки валидации)
	PrintStruct("Я просто строка, а не структура")
}