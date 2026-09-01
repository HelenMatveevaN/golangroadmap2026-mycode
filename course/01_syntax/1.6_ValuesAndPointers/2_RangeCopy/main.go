package main

import "fmt"

type Product struct {
	Name  string
	Price int
}

func main() {
	products := []Product{
		{Name: "Книга", Price: 500},
		{Name: "Ручка", Price: 50},
	}

	fmt.Println("--- Эксперимент 1: Пытаемся изменить данные через v ---")
	for _, v := range products {
		v.Price = 999 // Меняем цену во временной копии v!
	}
	// Смотрим на исходный слайс — цены НЕ изменились!
	fmt.Printf("После цикла с 'v': %v\n\n", products)

	fmt.Println("--- Эксперимент 2: Доказываем копирование через адреса памяти ---")
	for i, v := range products {
		// Адрес элемента внутри оригинального слайса
		origAddr := &products[i]
		// Адрес переменной v
		vAddr := &v

		fmt.Printf("Элемент [%d]: Адрес в слайсе: %p | Адрес переменной v: %p\n", i, origAddr, vAddr)
	}
}