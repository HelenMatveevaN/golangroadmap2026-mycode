package main

import "fmt"

type Product struct {
	Name  string
	Price int
}

func main() {
	// теперь мы создаем слайс УКАЗАТЕЛЕЙ []*Product	
	products := []*Product{
		{Name: "Книга", Price: 500},
		{Name: "Ручка", Price: 50},
	}

	fmt.Println("--- Эксперимент: Изменяем данные через копию указателя v ---")
	for _, v := range products {
		// v — это копия указателя. Она содержит адрес оригинальной структуры.
		// Go автоматически разадресовывает её под капотом, когда мы пишем v.Price
		v.Price = 999
	}

	// Проверяем исходный слайс — цены УСПЕШНО изменились!
	fmt.Println("После цикла со слайсом указателей:")
	for _, p := range products {
		fmt.Printf("Товар: %s, Цена: %d (Адрес объекта: %p)\n", p.Name, p.Price, p)
	}
}