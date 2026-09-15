package main

import (
	"fmt"
	"sync"
)

type Database struct {
	connection string
}

var (
	db   *Database
	once sync.Once
)

func GetDB() *Database {
	once.Do(func() {
		fmt.Println("Инициализируем соединение с БД...")
		db = &Database{connection: "postgres://localhost:5432/mydb"}
	})
	return db
}

func main() {
	var wg sync.WaitGroup

	// 10 горутин пытаются получить соединение
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			d := GetDB()
			fmt.Printf("Горутина %d использует: %s\n", id, d.connection)
		}(i)
	}

	wg.Wait()
	// "Инициализируем..." напечатается ровно один раз
}