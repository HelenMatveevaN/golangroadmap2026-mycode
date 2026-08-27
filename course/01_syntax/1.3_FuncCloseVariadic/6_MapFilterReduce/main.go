package main

//Реализовать map/filter/reduce через дженерики (или без них для тренировки).
//Функциональные примитивы (Map/Filter/Reduce)

import (
	"fmt"
	"strconv"
)

//[T, U any]: Это объявление плейсхолдеров (параметров) типов
//slice []T: Исходный слайс состоит из элементов первого типа T
//f func(T) U: Функция-трансформатор принимает на вход элемент типа T, а на выходе возвращает тип U
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))

	for i, item := range slice {
		result[i] = f(item) //функция трансформации
	}

	return result
}

// Filter фильтрует слайс на основе функции-предиката f.
// Возвращает новый слайс, содержащий только элементы, прошедшие проверку
func Filter[T any](slice []T, f func(T) bool)   []T {
	var result []T

	for _, item := range slice {
		if f(item) { //filter
			result = append(result, item)
		}
	}

	return result
}

//Сжимает слайс в одно единственное значение типа U 
// (например, сумму чисел или склеенную строку), 
// используя начальное значение init
func Reduce[T, U any](slice []T,  init U, f func(U, T) U) U {
	accumulator := init
	for _, item := range slice {
		accumulator = f(accumulator, item)
	}
	return accumulator
}


func main() {
	fmt.Println("=== Map ===")
	//Преобразование чисел (int -> int)
	numbers := []int{1,2,3,4,5}
	squares := Map(numbers, func(n int) int {
		return n * n
	})
	fmt.Printf("Квадраты чисел: %v\n", squares)

	//Трансформация типов (int -> string)
	stringNumbers := Map(numbers, func(n int) string {
		return "Число: " + strconv.Itoa(n)
	})
	fmt.Printf("Строковый слайс: %v\n", stringNumbers)

	//Работа со строками (string -> int)
	words := []string{"Go", "Generics", "2026"}
	lengths := Map(words, func(s string) int {
		return len(s)
	})
	fmt.Printf("Длины слов %#v: %v\n", words, lengths)


	fmt.Println("\n=== Filter ===")
	//Фильтрация чисел (int) — отбираем только четные
	numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evenNumbers := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("Четные числа в %v: %v\n", numbers, evenNumbers)

	//Фильтрация строк (string) — отбираем слова длиннее 3 символов
	words = []string{"Go", "Rust", "C", "Python", "PHP"}
	longWords := Filter(words, func(s string) bool {
		return len(s) > 3
	})
	fmt.Printf("Длинные слова в %v: %v\n", words, longWords)


	fmt.Println("\n=== Reduce ===")
	numbers = []int{10,20,30,40,50}

	//сжатие в 1 число(сумма)
	sum := Reduce(numbers, 0, func(acc int, item int) int {
		return acc + item
	})
	fmt.Printf("Сумма %v: %v\n", numbers, sum)

	// Сжатие в ОДНУ строчку (Текст)
	text := Reduce(words, "Результат склеивания строк:", func(acc string, item string) string {
		return fmt.Sprintf("%v %s", acc, item)
	})
	fmt.Println(text)
}