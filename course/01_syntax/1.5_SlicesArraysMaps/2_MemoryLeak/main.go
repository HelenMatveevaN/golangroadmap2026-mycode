package main

import (
	"fmt"
	"runtime"
	"time"
)

// Объявляем глобальные переменные на уровне пакета.
// Это гарантирует компилятору, что данные ИСПОЛЬЗУЮТСЯ глобально, 
// и заставит рантайм выделить память в Heap.
var GlobalLeaked []int
var GlobalFixed  []int

// LeakGenerator имитирует чтение огромного файла/лога
func LeakGenerator() []int {
	// Создаем массив на 10 миллионов интов (около 80 МБ в памяти)
	huge := make([]int, 10_000_000)
	huge[0] = 42
	huge[1] = 99

	// УТЕЧКА: Мы возвращаем срез всего из 2-х элементов.
	// Но под капотом возвращаемый Slice Header содержит ptr, 
	// который ссылается на этот огромный массив в 80 МБ!
	return huge[:2]
}

// FixGenerator делает то же самое, но безопасно
func FixGenerator() []int {
	// Создаем массив на 10 миллионов интов (около 80 МБ в памяти)
	huge := make([]int, 10_000_000)
	huge[0] = 42
	huge[1] = 99

	small := make([]int, 2)
	copy(small, huge[:2])

	return small
}

// Вспомогательная функция для отображения занятой памяти в Мегабайтах
func printMemUsage(msg string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%-25s Занято в Heap: %v MB\n", msg, m.Alloc / 1024 / 1024)
}

func main() {
	printMemUsage("1. Старт программы:")

	// --- ТЕСТ С УТЕЧКОЙ ---
	GlobalLeaked = LeakGenerator()

	// Принудительно вызываем сборщик мусора
	runtime.GC() 
	time.Sleep(50 * time.Millisecond) // даем GC время отработать

	printMemUsage("2. После LeakGenerator:")
	
	// Насильно очищаем глобальную ссылку, чтобы запустить второй тест с чистого листа
	GlobalLeaked = nil 
	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	// --- ТЕСТ С ИСПРАВЛЕНИЕМ ---
	GlobalFixed := FixGenerator()

	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	printMemUsage("3. После FixGenerator:")

	// Держим ссылку, просто чтобы компилятор не удалил GlobalFixed раньше времени
	_ = GlobalFixed

}