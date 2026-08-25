package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	// Строка содержит:
	// - Английские буквы (ASCII, 1 байт на символ)
	// - Кириллицу (UTF-8, 2 байта на символ)
	// - Эмодзи (UTF-8, 4 байта на символ)
	// 1 байт = 8 бит
	str := "GoПривет🚀"

	fmt.Println("=== ЧАСТЬ 1: Измерение длины ===")
	// Встроенная функция len() ВСЕГДА возвращает количество БАЙТ в строке
	fmt.Printf("len(str):               %d байт \n", len(str))

	// Функция из пакета unicode/utf8 считает честные символы (руны), декодируя UTF-8 на лету
	fmt.Printf("RuneCountInString(str):  %d рун \n", utf8.RuneCountInString(str))

	// Преобразование (каст) строки в срез рун создает новый массив в памяти
	fmt.Printf("len([]rune(str)):        %d элементов\n", len([]rune(str)))

	fmt.Println("\n=== ЧАСТЬ 2: Опасность индексации строки ===")
	// Попытка взять элемент строки по индексу возвращает БАЙТ (uint8), а не символ
	badChar := str[2]
	fmt.Printf("str[2] (первый байт буквы 'П'): в hex: 0x%x | как символ: %c\n", badChar, badChar)
	fmt.Println("Результат: мы сломали UTF-8 символ и получили 'битый' вывод.")

	fmt.Println("\n=== ЧАСТЬ 3: Побайтное представление ([]byte) ===")
	// Каст строки в []byte копирует данные в изменяемый срез
	byteSlice := []byte(str)
	fmt.Printf("Срез байт (hex): %x\n", byteSlice)
	fmt.Printf("Буква 'G' занимает 1 байт: 0x%x\n", byteSlice[0])
	fmt.Printf("Буква 'П' занимает 2 байта: 0x%x 0x%x\n", byteSlice[2], byteSlice[3])
	fmt.Printf("Эмодзи 🚀 занимает 4 байта: 0x%x 0x%x 0x%x 0x%x\n", byteSlice[14], byteSlice[15], byteSlice[16], byteSlice[17])

	fmt.Println("\n=== ЧАСТЬ 4: Порунное представление ([]rune) ===")
	// rune — это синоним int32. Каждый элемент гарантированно вмещает любой символ Юникода
	runeSlice := []rune(str)
	for i, r := range runeSlice {
		// %U выводит код Юникода (например, U+1F680)		
		fmt.Printf("Индекс [%d] -> Код: %-10U | Символ: %c\n", i, r, r)
	}

	fmt.Println("\n=== ЧАСТЬ 5: Иммутабельность строк ===")
	// str[0] = 'W' // ОШИБКА КОМПИЛЯЦИИ: cannot assign to str[0] (strings are immutable)

	// Чтобы изменить строку, её нужно превратить в срез, изменить его и вернуть обратно
	modifiedRunes := []rune(str)
	modifiedRunes[2] = 'п'
	newStr := string(modifiedRunes)
	fmt.Printf("Измененная строка: %s\n", newStr)
}
