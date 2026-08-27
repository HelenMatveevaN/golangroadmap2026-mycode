package main

//Замыкание-счётчик, замыкание-мемоизатор.
//Эволюция замыканий (Счётчик и Мемоизатор)

import (
	"fmt"
	"time"
)

func NewCounter() func() int {
	count :=0 //escape to heap

	return func() int {
		count++
		return count
	}
}

// Memoize принимает функцию и возвращает её кэшированную версию
func Memoize(f func(int) int) func(int) int {
	// Карта создается один раз при инициализации кэша и убегает в кучу
	cache := make(map[int]int)

	// Возвращаем анонимную функцию-обертку
	return func(x int) int {
		if value, ok := cache[x]; ok {
			return value
		}

		//если в кэше пусто
		result := f(x) //вызываем тяжелую функцию
		cache[x] = result
		return result
	}
}

//тяжелая функция
func slowSquare(n int) int {
	fmt.Println("--> (Тяжелое вычисление для числа", n, ")")
	time.Sleep(2 * time.Second)
	return n * n
}

func main() {
	//Counter
	fmt.Println("=== Counter ===")
	genCnt := NewCounter()
	fmt.Println(genCnt())
	fmt.Println(genCnt())

	genCnt2 := NewCounter()
	fmt.Println(genCnt2())
	fmt.Println(genCnt2())
	fmt.Println(genCnt2())

	//Memoize
	fmt.Println("=== Memoize ===")	
	cachedSquare := Memoize(slowSquare)

	fmt.Println("Считаем квадрат числа 5...")
	start := time.Now()
	fmt.Println("Результат:", cachedSquare(5)) //2sec
	fmt.Printf("Время выполнения: %v\n\n", time.Since(start))

	fmt.Println("Снова просим квадрат числа 5...")
	start = time.Now()
	fmt.Println("Результат:", cachedSquare(5))
	fmt.Printf("Время выполнения: %v\n\n", time.Since(start))

	fmt.Println("Считаем квадрат числа 10...")
	start = time.Now()
	fmt.Println("Результат:", cachedSquare(10)) //2sec
	fmt.Printf("Время выполнения: %v\n\n", time.Since(start))	
}